/* OpenStories hotsite. No framework. The stress-tester runs against the live API when
   window.OPENSTORIES_API is set; otherwise it mirrors internal/evaluator/evaluator.go
   in the browser over a compact copy of the library (data/stories.min.json). */
(() => {
  const $ = (s, r = document) => r.querySelector(s);
  const $$ = (s, r = document) => [...r.querySelectorAll(s)];
  const API = (window.OPENSTORIES_API || "").replace(/\/$/, "");
  const REPO = "https://github.com/gabrielrondon/openstories";

  // ---- counters
  const easeOut = (t) => 1 - Math.pow(1 - t, 3);
  const reduce = matchMedia("(prefers-reduced-motion: reduce)").matches;
  $$("[data-count]").forEach((el) => {
    const target = +el.dataset.count;
    if (reduce) { el.textContent = target.toLocaleString(); return; }
    const t0 = performance.now();
    const tick = (now) => {
      const p = Math.min(1, (now - t0) / 1400);
      el.textContent = Math.round(target * easeOut(p)).toLocaleString();
      if (p < 1) requestAnimationFrame(tick);
    };
    requestAnimationFrame(tick);
  });

  // ---- copy buttons
  $$(".copy").forEach((b) => b.addEventListener("click", async () => {
    try { await navigator.clipboard.writeText(b.dataset.copy); b.classList.add("done"); b.querySelector(".hint").textContent = "copied"; setTimeout(() => { b.classList.remove("done"); b.querySelector(".hint").textContent = "copy"; }, 1600); } catch { /* no clipboard */ }
  }));

  // ---- data
  const load = (p) => fetch(p).then((r) => r.json());
  let corpus = null;
  const corpusP = () => corpus || (corpus = load("/data/stories.min.json"));

  // ---- evaluator mirror (internal/evaluator/evaluator.go + internal/store/store.go)
  // extractKeyPhrases: ASCII words longer than 4 chars.
  const phrases = (s) => s.toLowerCase().split(/[^a-z0-9]+/).filter((w) => w.length > 4);
  // store.tokenize: words longer than 1 char, minus stop words (keeps à-ú range like Go).
  const STOP = new Set(["a","o","de","do","da","em","um","uma","para","com","e","ou","the","and","of","to","in","for","is","on","that","this","as","i","want","so"]);
  const tokenize = (s) => s.toLowerCase().split(/[^a-z0-9\u00e0-\u00fa]+/).filter((w) => w.length > 1 && !STOP.has(w));
  // store.calculateScore: weighted hits, then demand multiplier.
  function calculateScore(s, toks) {
    let score = 0;
    const title = s.t.toLowerCase(), asA = s.a.toLowerCase(), iWant = s.w.toLowerCase(), soThat = s.o.toLowerCase();
    const ecs = s.e.map((x) => x.toLowerCase()), qs = s.q.map((x) => x.toLowerCase()), acs = s.c.map((x) => x.toLowerCase());
    for (const t of toks) {
      if (title.includes(t)) score += 25;
      if (iWant.includes(t)) score += 15;
      if (soThat.includes(t) || asA.includes(t)) score += 10;
      for (const g of s.g) if (g.toLowerCase() === t) score += 15;
      for (const ec of ecs) if (ec.includes(t)) score += 8;
      for (const q of qs) if (q.includes(t)) score += 6;
      for (const ac of acs) if (ac.includes(t)) score += 6;
    }
    return score > 0 ? score * (1 + s.s / 10) : 0;
  }
  // store.Search: filter by industry/domain, rank by calculateScore (or demand when no query).
  function search(stories, q, industry, domain) {
    const toks = tokenize(q);
    const scored = [];
    for (const s of stories) {
      if (industry && s.n.toLowerCase() !== industry.toLowerCase()) continue;
      if (domain && s.d.toLowerCase() !== domain.toLowerCase()) continue;
      if (!toks.length) { scored.push({ s, k: s.s }); continue; }
      const k = calculateScore(s, toks);
      if (k > 0) scored.push({ s, k });
    }
    scored.sort((a, b) => b.k - a.k);
    return scored.map((x) => x.s);
  }

  function evaluateLocal(stories, spec, industry, domain) {
    let matches = search(stories, spec, industry, domain);
    if (!matches.length && industry) matches = search(stories, "", industry, domain);
    const specLower = spec.toLowerCase();
    let total = 0, covered = 0;
    const matched = [], alerts = [];
    for (const st of matches.slice(0, 5)) {
      let ok = true; const missed = [];
      for (const ec of st.e) {
        total++;
        const hit = phrases(ec).some((t) => specLower.includes(t));
        if (hit) covered++; else { ok = false; missed.push(ec); alerts.push({ id: st.i, title: st.t, ec }); }
      }
      matched.push({ story_id: st.i, title: st.t, demand_score: st.s, covered: ok, missed_cases: missed });
    }
    let score, summary;
    if (!total) { score = 75; summary = "No domain-specific stories with failure criteria found for detailed comparison. Specification appears reasonable but lacks empirical validation."; }
    else {
      score = Math.floor((covered / total) * 100);
      summary = score >= 85 ? "High resilience! The technical specification preemptively mitigates known production failure modes documented in field reports."
        : score >= 60 ? "Moderate coverage. The happy path is addressed, but critical production blind spots frequently reported by users in the wild are missing."
        : "Severe Production Vulnerability Alert: Multiple high-impact failure modes observed in real-world outages were ignored in this specification.";
    }
    return { score, summary, matched_stories: matched, edge_case_alerts: alerts.map((a) => `[${a.id}] ${a.title}: ${a.ec}`) };
  }
  window.OpenStoriesEval = { tokenize, search, evaluateLocal, corpus: corpusP };
  async function evaluate(spec, industry, domain) {
    if (API) {
      const r = await fetch(`${API}/api/eval`, { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify({ spec, industry, domain }) });
      if (r.ok) return { result: await r.json(), mode: "live API" };
    }
    return { result: evaluateLocal(await corpusP(), spec, industry, domain), mode: "in your browser, same rules as openstories eval" };
  }

  // ---- render report
  const report = $("#report"), gauge = $("#gaugeFill"), scoreNum = $("#scoreNum"), summary = $("#summary"), alerts = $("#alerts"), noAlerts = $("#noAlerts"), mode = $("#mode");
  const storyLink = (id) => `${REPO}/search?q=${encodeURIComponent(id)}&type=code`;
  function render(res, how) {
    report.hidden = false;
    const s = Math.max(0, Math.min(100, res.score));
    report.dataset.band = s >= 85 ? "high" : s >= 60 ? "mid" : "low";
    scoreNum.textContent = s;
    requestAnimationFrame(() => { gauge.style.strokeDashoffset = String(326.7 * (1 - s / 100)); });
    summary.textContent = res.summary;
    alerts.innerHTML = "";
    const list = (res.edge_case_alerts || []).slice(0, 8);
    for (const a of list) {
      const m = a.match(/^\[([^\]]+)\]\s*(.*?):\s*(.*)$/);
      const li = document.createElement("li");
      if (m) li.innerHTML = `<span><b><a href="${storyLink(m[1])}" rel="noopener">${m[1]}</a></b> ${escapeHtml(m[3])}<br><small>${escapeHtml(m[2])}</small></span>`;
      else li.innerHTML = `<span>${escapeHtml(a)}</span>`;
      alerts.appendChild(li);
    }
    noAlerts.hidden = list.length > 0;
    mode.textContent = how;
    report.scrollIntoView({ block: "nearest", behavior: reduce ? "auto" : "smooth" });
  }
  const escapeHtml = (s) => s.replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c]));

  // ---- samples and form
  const spec = $("#spec"), industry = $("#industry"), samplesEl = $("#samples");
  load("/data/taxonomies.json").then((tax) => {
    industry.innerHTML = `<option value="">any industry</option>` + tax.map((t) => `<option value="${t.industry}">${t.industry}</option>`).join("");
  });
  load("/data/playground.json").then((samples) => {
    samples.forEach((sm, i) => {
      const b = document.createElement("button");
      b.type = "button"; b.setAttribute("role", "listitem");
      b.textContent = sm.spec.length > 64 ? sm.spec.slice(0, 62) + "…" : sm.spec;
      b.title = sm.spec;
      b.addEventListener("click", () => {
        $$("button", samplesEl).forEach((x) => x.setAttribute("aria-pressed", "false"));
        b.setAttribute("aria-pressed", "true");
        spec.value = sm.spec; industry.value = sm.industry;
        render(sm.result, "precomputed with openstories eval v1.0");
      });
      samplesEl.appendChild(b);
      if (i === 0) { spec.value = sm.spec; industry.value = sm.industry; }
    });
  });
  $("#specForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const text = spec.value.trim(); if (!text) return;
    $$("button", samplesEl).forEach((x) => x.setAttribute("aria-pressed", "false"));
    mode.textContent = "testing…";
    const { result, mode: how } = await evaluate(text, industry.value, "");
    render(result, how);
  });

  // ---- incident showcase
  load("/data/showcase.json").then((stories) => {
    const wrap = $("#incidentCards");
    for (const s of stories) {
      const ac = s.acceptance_criteria?.[0];
      const ev = s.evidence?.[0];
      const card = document.createElement("article");
      card.className = "card";
      card.innerHTML = `
        <div class="meta"><span>${s.id}</span><span>${s.industry} / ${s.domain}</span><span class="score">demand ${s.demand_score}</span></div>
        <h3>${escapeHtml(s.title)}</h3>
        <p class="story"><b>As a</b> ${escapeHtml(s.story.as_a)}, <b>I want</b> ${escapeHtml(s.story.i_want)}, <b>so that</b> ${escapeHtml(s.story.so_that)}</p>
        ${ac ? `<pre class="gherkin"><b>Scenario:</b> ${escapeHtml(ac.scenario)}\n<b>Given</b> ${escapeHtml(ac.given)}\n<b>When</b> ${escapeHtml(ac.when)}\n<b>Then</b> ${escapeHtml(ac.then)}</pre>` : ""}
        ${ev ? `<blockquote>“${escapeHtml(ev.quote)}”<a href="${ev.source}" rel="noopener">${escapeHtml(ev.source.replace(/^https?:\/\//, ""))}${ev.date ? " · " + escapeHtml(ev.date) : ""}</a></blockquote>` : ""}
      `;
      wrap.appendChild(card);
    }
  });
})();
