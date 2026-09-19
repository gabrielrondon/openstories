#!/usr/bin/env python3
"""Build site/data/*.json from a running `openstories serve` instance.

Usage: python3 scripts/build-data.py [http://127.0.0.1:8090]

Writes:
  data/stories.min.json   compact corpus with exactly the fields the Go store
                          scores on (store.calculateScore) and the evaluator
                          checks (evaluator.Evaluate), in store order.
  data/taxonomies.json    industries and domains.
  data/playground.json    precomputed evals for the sample specs (same API).
"""
import json, sys, urllib.request

API = (sys.argv[1] if len(sys.argv) > 1 else "http://127.0.0.1:8090").rstrip("/")
HERE = __import__("os").path.dirname(__import__("os").path.abspath(__file__))
DATA = __import__("os").path.join(HERE, "..", "data")

def get(path, body=None):
    req = urllib.request.Request(API + path, data=json.dumps(body).encode() if body else None,
                                 headers={"content-type": "application/json"} if body else {})
    with urllib.request.urlopen(req) as r:
        return json.load(r)

stories = get("/api/stories")
out = []
for s in stories:
    st = s.get("story") or {}
    out.append({
        "i": s["id"], "n": s["industry"], "d": s["domain"], "t": s["title"],
        "s": s.get("demand_score", 0),
        "a": st.get("as_a", ""), "w": st.get("i_want", ""), "o": st.get("so_that", ""),
        "g": s.get("tags") or [],
        "e": s.get("edge_cases") or [],
        "q": [ev.get("quote", "") for ev in (s.get("evidence") or [])],
        "c": [" ".join([ac.get("scenario", ""), ac.get("given", ""), ac.get("when", ""), ac.get("then", "")])
              for ac in (s.get("acceptance_criteria") or [])],
    })
with open(f"{DATA}/stories.min.json", "w") as f:
    json.dump(out, f, separators=(",", ":"), ensure_ascii=False)

with open(f"{DATA}/taxonomies.json", "w") as f:
    json.dump(get("/api/taxonomies"), f, separators=(",", ":"), ensure_ascii=False)

pg_path = f"{DATA}/playground.json"
samples = json.load(open(pg_path))
for smp in samples:
    smp["result"] = get("/api/eval", {"spec": smp["spec"], "industry": smp["industry"], "domain": smp.get("domain", "")})
with open(pg_path, "w") as f:
    json.dump(samples, f, indent=1, ensure_ascii=False)

print(f"stories: {len(out)}  playground samples: {len(samples)}")
