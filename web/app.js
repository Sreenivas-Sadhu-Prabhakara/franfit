/* FranFit frontend — vanilla JS, no dependencies. */
(() => {
  "use strict";

  const $ = (sel) => document.querySelector(sel);
  const $$ = (sel) => Array.from(document.querySelectorAll(sel));

  // ---------- INR formatting (mirrors internal/money) ----------
  function formatLakhs(l) {
    if (l >= 100) return "₹" + trim1(l / 100) + " Cr";
    if (l >= 1) return "₹" + trim1(l) + " L";
    return formatRupees(Math.round(l * 100000));
  }
  function formatRupees(n) {
    const sign = n < 0 ? "-" : "";
    let s = String(Math.abs(n));
    if (s.length <= 3) return sign + "₹" + s;
    let head = s.slice(0, -3), tail = s.slice(-3), groups = [];
    while (head.length > 2) { groups.unshift(head.slice(-2)); head = head.slice(0, -2); }
    if (head) groups.unshift(head);
    return sign + "₹" + groups.join(",") + "," + tail;
  }
  const trim1 = (f) => String(Math.round(f * 10) / 10);
  const esc = (s) => String(s).replace(/[&<>"']/g, (c) =>
    ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c]));

  // ---------- state ----------
  const state = {
    step: 1,
    budgetL: 30,
    city: "",
    tier: 2,
    involvement: "owner-operator",
    risk: 3,
    categories: [],
    spaceSqft: 0,
    matches: [],
    activeBrand: null,   // brand shown in drawer
    activeMatch: null,   // match backing the drawer / lead form
  };
  const riskLabels = { 1: "Safety first", 2: "Cautious", 3: "Balanced", 4: "Growth-hungry", 5: "Aggressive" };
  let categories = ["QSR", "Cafe", "Salon", "Pharmacy", "Education", "Grocery", "Fitness", "Courier"];

  // ---------- navigation ----------
  function show(view) {
    $("#view-discover").hidden = view !== "discover";
    $("#view-leads").hidden = view !== "leads";
    $$(".nav-link").forEach((b) => b.classList.toggle("is-active", b.dataset.nav === view));
    if (view === "leads") loadLeads();
  }
  document.addEventListener("click", (e) => {
    const nav = e.target.closest("[data-nav]");
    if (nav) { e.preventDefault(); show(nav.dataset.nav); }
  });

  // ---------- wizard ----------
  function setStep(n) {
    state.step = n;
    $$(".wiz-step").forEach((el) => el.classList.toggle("is-active", +el.dataset.step === n));
    $$(".wiz-tab").forEach((el) => {
      el.classList.toggle("is-active", +el.dataset.step === n);
      el.classList.toggle("is-done", +el.dataset.step < n);
    });
    $("#btnBack").hidden = n === 1;
    $("#btnNext").textContent = n === 3 ? "Find my matches" : "Next";
  }
  $("#btnBack").addEventListener("click", () => setStep(state.step - 1));
  $("#btnNext").addEventListener("click", () => {
    if (state.step < 3) { setStep(state.step + 1); return; }
    runMatch();
  });
  $$(".wiz-tab").forEach((tab) =>
    tab.addEventListener("click", () => setStep(+tab.dataset.step)));

  // Step 1 inputs
  const budgetRange = $("#budgetRange");
  const budgetOut = $("#budgetReadout");
  budgetRange.addEventListener("input", () => {
    state.budgetL = +budgetRange.value;
    budgetOut.textContent = formatLakhs(state.budgetL);
  });
  $("#city").addEventListener("input", (e) => { state.city = e.target.value.trim(); });
  $("#tierChips").addEventListener("click", (e) => {
    const chip = e.target.closest(".chip");
    if (!chip) return;
    state.tier = +chip.dataset.tier;
    $$("#tierChips .chip").forEach((c) => {
      const on = c === chip;
      c.classList.toggle("is-on", on);
      c.setAttribute("aria-checked", String(on));
    });
  });
  $("#space").addEventListener("input", (e) => { state.spaceSqft = Math.max(0, +e.target.value || 0); });

  // Step 2 inputs
  $("#invCards").addEventListener("click", (e) => {
    const card = e.target.closest(".inv-card");
    if (!card) return;
    state.involvement = card.dataset.inv;
    $$("#invCards .inv-card").forEach((c) => {
      const on = c === card;
      c.classList.toggle("is-on", on);
      c.setAttribute("aria-checked", String(on));
    });
  });
  const riskInput = $("#risk");
  riskInput.addEventListener("input", () => {
    state.risk = +riskInput.value;
    $("#riskReadout").innerHTML = `<span class="num">${state.risk}</span>/5 — ${riskLabels[state.risk]}`;
  });

  // Step 3: category chips
  function renderCategoryChips() {
    $("#catChips").innerHTML = categories.map((c) =>
      `<button class="chip" data-cat="${esc(c)}">${esc(c)}</button>`).join("");
  }
  $("#catChips").addEventListener("click", (e) => {
    const chip = e.target.closest(".chip");
    if (!chip) return;
    const cat = chip.dataset.cat;
    const i = state.categories.indexOf(cat);
    if (i >= 0) state.categories.splice(i, 1); else state.categories.push(cat);
    chip.classList.toggle("is-on", i < 0);
  });

  // ---------- matching ----------
  async function runMatch() {
    const btn = $("#btnNext");
    btn.disabled = true;
    btn.textContent = "Scoring 30 brands…";
    try {
      const res = await fetch("/api/v1/match", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          budgetL: state.budgetL, city: state.city, tier: state.tier,
          involvement: state.involvement, risk: state.risk,
          categories: state.categories, spaceSqft: state.spaceSqft,
        }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "match failed");
      state.matches = data.matches || [];
      renderResults(data);
    } catch (err) {
      toast("Could not score matches: " + err.message);
    } finally {
      btn.disabled = false;
      btn.textContent = "Find my matches";
    }
  }

  function dialSVG(score) {
    const r = 30, c = 2 * Math.PI * r, filled = c * score / 100;
    const color = score >= 75 ? "var(--leaf)" : score >= 50 ? "var(--marigold-deep)" : "var(--madder)";
    let ticks = "";
    for (let i = 0; i < 24; i++) {
      const a = (i / 24) * 2 * Math.PI - Math.PI / 2;
      const x1 = 39 + Math.cos(a) * 36, y1 = 39 + Math.sin(a) * 36;
      const x2 = 39 + Math.cos(a) * 38, y2 = 39 + Math.sin(a) * 38;
      ticks += `<line x1="${x1.toFixed(1)}" y1="${y1.toFixed(1)}" x2="${x2.toFixed(1)}" y2="${y2.toFixed(1)}" stroke="var(--line)" stroke-width="1.4"/>`;
    }
    return `<svg class="dial" viewBox="0 0 78 78" role="img" aria-label="Fit score ${score} out of 100">
      ${ticks}
      <circle cx="39" cy="39" r="${r}" fill="none" stroke="rgba(18,49,46,.08)" stroke-width="6"/>
      <circle cx="39" cy="39" r="${r}" fill="none" stroke="${color}" stroke-width="6" stroke-linecap="round"
        stroke-dasharray="${filled.toFixed(1)} ${c.toFixed(1)}" transform="rotate(-90 39 39)"/>
      <text x="39" y="42" text-anchor="middle" font-size="19" font-weight="700" fill="var(--ink)">${score}</text>
      <text x="39" y="54" text-anchor="middle" font-size="8.5" fill="var(--ink-soft)">/100</text>
    </svg>`;
  }

  function renderResults(data) {
    const results = $("#results");
    results.hidden = false;
    const cityBit = state.city ? ` in ${esc(state.city)}` : "";
    if (data.noMatches) {
      $("#matchList").innerHTML = "";
      $("#resultsSub").textContent = "";
      $("#noMatch").hidden = false;
      $("#noMatchText").textContent = data.explanation;
      $("#resultsTitle").textContent = "Your matches";
    } else {
      $("#noMatch").hidden = true;
      $("#resultsTitle").textContent = `Your top ${Math.min(state.matches.length, 30)} matches`;
      $("#resultsSub").textContent =
        `${state.matches.length} of 30 brands fit your ${formatLakhs(state.budgetL)} budget${cityBit} — ranked by fit score.`;
      $("#matchList").innerHTML = state.matches.map(matchCard).join("");
    }
    results.scrollIntoView({ behavior: "smooth", block: "start" });
  }

  function matchCard(m, i) {
    const b = m.brand;
    const factors = m.factors.map((f) => {
      const pct = f.max ? (f.points / f.max) * 100 : 0;
      const cls = pct >= 75 ? "" : pct >= 40 ? "mid" : "low";
      return `<div class="factor">
        <span class="factor-label">${esc(f.label)}</span>
        <span class="factor-bar"><span class="factor-fill ${cls}" style="width:${pct.toFixed(0)}%"></span></span>
        <span class="factor-pts">${f.points}/${f.max}</span>
        <span class="factor-note">${esc(f.note)}</span>
      </div>`;
    }).join("");
    return `<article class="match card" data-idx="${i}" style="animation-delay:${Math.min(i * 45, 400)}ms">
      <div class="match-top">
        ${dialSVG(m.score)}
        <div class="match-info">
          <h3 class="match-name">${esc(b.name)}</h3>
          <div class="match-meta">
            <span class="tag">${esc(b.category)}</span>
            <span class="tag model">${esc(m.recommendedModel)} for you</span>
            <span class="num">${formatLakhs(b.investmentMinL)}–${formatLakhs(b.investmentMaxL)}</span>
            <span>payback ≈ <span class="num">${b.paybackMonthsEst} mo</span></span>
          </div>
          <button class="expand-toggle" data-act="expand">Why this score ▾</button>
        </div>
        <div class="match-actions">
          <button class="btn ghost small" data-act="details">Details</button>
          <button class="btn primary small" data-act="intro">Request intro</button>
        </div>
      </div>
      <div class="match-expand">
        ${factors}
        <p class="match-reason">${esc(m.reasoning)}</p>
      </div>
    </article>`;
  }

  $("#matchList").addEventListener("click", (e) => {
    const btn = e.target.closest("[data-act]");
    if (!btn) return;
    const card = btn.closest(".match");
    const m = state.matches[+card.dataset.idx];
    if (btn.dataset.act === "expand") {
      card.classList.toggle("is-open");
      btn.textContent = card.classList.contains("is-open") ? "Hide breakdown ▴" : "Why this score ▾";
    } else if (btn.dataset.act === "details") {
      openDrawer(m.brand, m);
    } else if (btn.dataset.act === "intro") {
      openLeadForm(m.brand, m);
    }
  });

  $("#btnRedo").addEventListener("click", () => {
    $("#results").hidden = true;
    setStep(1);
    $("#wizard").scrollIntoView({ behavior: "smooth" });
  });
  $("#btnRaise").addEventListener("click", () => {
    $("#results").hidden = true;
    setStep(1);
    budgetRange.focus();
    $("#wizard").scrollIntoView({ behavior: "smooth" });
  });

  // ---------- drawer ----------
  function openDrawer(brand, match) {
    state.activeBrand = brand;
    state.activeMatch = match || null;
    $("#dCat").textContent = brand.category;
    $("#dName").textContent = brand.name;
    $("#dStory").textContent = brand.brandStory;
    $("#dInvest").textContent = `${formatLakhs(brand.investmentMinL)} – ${formatLakhs(brand.investmentMaxL)}`;
    $("#dFee").textContent = formatLakhs(brand.franchiseFeeL);
    $("#dRoyalty").textContent = brand.royaltyPct + "%";
    $("#dArea").textContent = brand.areaSqftMin.toLocaleString("en-IN") + " sqft";
    $("#dPayback").textContent = brand.paybackMonthsEst + " months";
    $("#dProfit").textContent = formatLakhs(brand.monthlyProfitEstL) + " / month";
    $("#dModels").textContent = brand.modelsSupported.join(" · ");
    $("#dTiers").textContent = brand.cityTiers.map((t) => "Tier " + t).join(" · ");
    const reason = $("#dReason");
    if (match) { reason.hidden = false; $("#dReasonText").textContent = match.reasoning; }
    else reason.hidden = true;
    $("#drawerScrim").hidden = false;
    const d = $("#drawer");
    d.setAttribute("aria-hidden", "false");
    requestAnimationFrame(() => d.classList.add("is-open"));
  }
  function closeDrawer() {
    const d = $("#drawer");
    d.classList.remove("is-open");
    d.setAttribute("aria-hidden", "true");
    $("#drawerScrim").hidden = true;
  }
  $("#drawerClose").addEventListener("click", closeDrawer);
  $("#drawerScrim").addEventListener("click", closeDrawer);
  $("#drawerIntro").addEventListener("click", () => {
    closeDrawer();
    openLeadForm(state.activeBrand, state.activeMatch);
  });

  // ---------- lead form ----------
  function openLeadForm(brand, match) {
    state.activeBrand = brand;
    state.activeMatch = match || null;
    $("#leadBrandName").textContent = brand.name;
    $("#leadBudget").value = state.budgetL;
    $("#leadError").hidden = true;
    $("#modalScrim").hidden = false;
    $("#leadName").focus();
  }
  function closeLeadForm() { $("#modalScrim").hidden = true; }
  $("#leadCancel").addEventListener("click", closeLeadForm);
  $("#modalScrim").addEventListener("click", (e) => { if (e.target === e.currentTarget) closeLeadForm(); });
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") { closeLeadForm(); closeDrawer(); }
  });

  $("#leadForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const btn = $("#leadSubmit");
    btn.disabled = true;
    try {
      const res = await fetch("/api/v1/leads", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          brandId: state.activeBrand.id,
          name: $("#leadName").value.trim(),
          phone: $("#leadPhone").value.trim(),
          email: $("#leadEmail").value.trim(),
          budgetL: +$("#leadBudget").value,
          fitScore: state.activeMatch ? state.activeMatch.score : 0,
          city: state.city,
        }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "could not save lead");
      closeLeadForm();
      $("#leadForm").reset();
      toast(`Intro requested — ${data.brandName}'s team will call. Ref ${data.id}`);
      refreshLeadCount();
    } catch (err) {
      const el = $("#leadError");
      el.textContent = err.message;
      el.hidden = false;
    } finally {
      btn.disabled = false;
    }
  });

  // ---------- leads view ----------
  async function loadLeads() {
    const res = await fetch("/api/v1/leads");
    const data = await res.json();
    const leads = data.leads || [];
    $("#leadsEmpty").hidden = leads.length > 0;
    $("#leadsTable").style.display = leads.length ? "" : "none";
    $("#leadsBody").innerHTML = leads.map((l) => {
      const cls = l.fitScore >= 75 ? "" : l.fitScore >= 50 ? "mid" : "low";
      const when = new Date(l.createdAt).toLocaleString("en-IN", {
        day: "2-digit", month: "short", hour: "2-digit", minute: "2-digit",
      });
      return `<tr>
        <td><div class="lead-name">${esc(l.name)}</div><div class="lead-id">${esc(l.id)}</div></td>
        <td class="lead-contact">${esc(l.phone)}<br>${esc(l.email)}</td>
        <td>${esc(l.city || "—")}</td>
        <td class="num">${formatLakhs(l.budgetL)}</td>
        <td>${esc(l.brandName)}</td>
        <td><span class="fit-pill ${cls}">${l.fitScore}</span></td>
        <td class="num">${esc(when)}</td>
      </tr>`;
    }).join("");
    updateLeadCount(leads.length);
  }
  function updateLeadCount(n) {
    const el = $("#leadCount");
    el.hidden = !n;
    el.textContent = n;
  }
  async function refreshLeadCount() {
    try {
      const res = await fetch("/api/v1/leads");
      const data = await res.json();
      updateLeadCount((data.leads || []).length);
    } catch { /* non-fatal */ }
  }

  // ---------- toast ----------
  let toastTimer;
  function toast(msg) {
    const el = $("#toast");
    el.textContent = msg;
    el.hidden = false;
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => { el.hidden = true; }, 4200);
  }

  // ---------- boot ----------
  async function boot() {
    renderCategoryChips();
    try {
      const res = await fetch("/api/v1/brands");
      const data = await res.json();
      if (data.categories && data.categories.length) {
        categories = data.categories;
        renderCategoryChips();
      }
      $("#heroBrandCount").textContent = data.count;
    } catch { /* directory is embedded; UI still works */ }
    refreshLeadCount();
  }
  boot();
})();
