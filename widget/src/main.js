import createGlobe from "cobe";

const STATUSES = ["resolved", "stale", "pending", "error", "nxdomain"];

const STATUS_RGB = {
  resolved: [0.09, 0.64, 0.29],
  stale: [0.92, 0.55, 0.05],
  pending: [0.45, 0.45, 0.5],
  error: [166, 25, 46],
  nxdomain: [246, 130, 31]
};

const STATUS_TEXT = {
  resolved: "Resolved",
  stale: "Stale",
  pending: "Pending",
  error: "Error",
  nxdomain: "Non Existant"
};

const $ = (id) => document.getElementById(id);
const globeEl = $("globe");
const canvas = $("cobe");

const el = (tag, className, text) => {
  const node = document.createElement(tag);
  if (className) node.className = className;
  if (text != null) node.textContent = text;
  return node;
};

const anchorTo = (node, anchor, visibleVar) => {
  node.style.positionAnchor = anchor;
  node.style.opacity = `var(${visibleVar}, 0)`;
  node.style.filter = `blur(calc((1 - var(${visibleVar}, 0)) * 8px))`;
};

const toRad = (deg) => (deg * Math.PI) / 180;

const dpr = Math.min(window.devicePixelRatio || 1, 2);
let size = globeEl.offsetWidth;
let phi = 0;
let theta = 0.25;
let velocityPhi = 0;
let velocityTheta = 0;
let pointer = null;

const globe = createGlobe(canvas, {
  devicePixelRatio: dpr,
  width: size * dpr,
  height: size * dpr,
  phi,
  theta,
  dark: 0,
  diffuse: 1.2,
  mapSamples: 16000,
  mapBrightness: 6,
  mapBaseBrightness: 0,
  baseColor: [1, 1, 1],
  markerColor: [0.1, 0.1, 0.1],
  glowColor: [1, 1, 1],
  markerElevation: 0.02,
  markers: [],
});

new ResizeObserver(() => {
  size = globeEl.offsetWidth;
  globe.update({ width: size * dpr, height: size * dpr });
}).observe(globeEl);

canvas.addEventListener("pointerdown", (e) => {
  pointer = { x: e.clientX, y: e.clientY };
  velocityPhi = velocityTheta = 0;
  canvas.setPointerCapture(e.pointerId);
});
canvas.addEventListener("pointermove", (e) => {
  if (!pointer) return;
  velocityPhi = ((e.clientX - pointer.x) / size) * 3;
  velocityTheta = ((e.clientY - pointer.y) / size) * 3;
  pointer = { x: e.clientX, y: e.clientY };
  phi += velocityPhi;
  theta = Math.max(-1.2, Math.min(1.2, theta + velocityTheta));
});
const endDrag = () => (pointer = null);
canvas.addEventListener("pointerup", endDrag);
canvas.addEventListener("pointercancel", endDrag);

(function animate() {
  if (!pointer) {
    velocityPhi *= 0.95;
    velocityTheta *= 0.9;
    phi += 0.003 + velocityPhi;
    theta += velocityTheta + (0.25 - theta) * 0.02;
  }
  globe.update({ phi, theta });
  requestAnimationFrame(animate);
})();

requestAnimationFrame(() => canvas.classList.add("ready"));

let overlays = [];

function renderGlobe(results) {
  overlays.forEach((n) => n.remove());
  overlays = [];

  for (const r of results) {
    const node = el("div", `site site--${r.status}`);
    const pyramid = el("div", "pyramid");
    for (let i = 0; i < 4; i++) pyramid.append(el("div", "pyramid-face"));
    const chip = el("span", "site-chip");
    chip.append(el("span", "dot"), el("span", "site-name", r.resolver.toLowerCase()), el("span", "site-status", r.status));
    node.append(pyramid, chip);
    anchorTo(node, `--cobe-${r.id}`, `--cobe-visible-${r.id}`);
    overlays.push(node);
  }

  globeEl.append(...overlays);

  globe.update({
    markers: results.map((r) => ({ id: r.id, location: [r.lat, r.lng], size: 0.03, color: STATUS_RGB[r.status] })),
  });
}

function renderPanel(results) {
  const counts = Object.fromEntries(STATUSES.map((s) => [s, results.filter((r) => r.status === s).length]));
//   const total = results.length || 1; // TODO: We will take this later

  const progress = $("progress");
  progress.replaceChildren(
    ...STATUSES.filter((s) => counts[s]).map((s) => {
      const seg = el("span", `seg seg--${s}`);
      seg.style.flexGrow = counts[s];
      return seg;
    }),
  );

  $("summary").replaceChildren(
    ...STATUSES.map((s) => {
      const li = el("li", `badge badge--${s}`);
      li.append(el("span", "dot"), el("span", "badge-count", String(counts[s])), el("span", null, STATUS_TEXT[s]));
      return li;
    }),
  );
}

function focus(results) {
  const target = results.find((r) => r.status !== "resolved") ?? results[0];
  if (!target) return;
  phi = Math.PI - (toRad(target.lng) - Math.PI / 2);
}

function setResults(resolverResults) {
  const results = resolverResults.map((r, i) => ({ ...r, id: `resolver-${i}` }));
  renderGlobe(results);
  renderPanel(results);
  focus(results);
}

function setInput({ domain, recordType }) {
  if (typeof domain !== "string" || !domain) return;
  const title = $("domain");
  title.replaceChildren(document.createTextNode(domain.slice(0, 253)));
  if (typeof recordType === "string" && recordType) title.append(el("span", "record-type", recordType.slice(0, 10)));
}

function render(sample) {
  setInput(sample);
  setResults(sample.results);
}

window.updateGlobe = render;

render(window.__DNS_DATA__ ?? []);
