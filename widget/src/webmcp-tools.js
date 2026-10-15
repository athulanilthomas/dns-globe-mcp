export function registerTools(tools) {
  const ctx = document.modelContext ?? navigator.modelContext;
  if (typeof ctx?.registerTool !== "function") {
    console.warn("[dns-globe] WebMCP is unavailable (no document.modelContext); tools not registered.");
    return;
  }

  for (const tool of tools) {
    const fail = (err) => console.warn(`[dns-globe] registerTool("${tool.name}") failed (duplicate after HMR?):`, err);
    try {
      Promise.resolve(ctx.registerTool(tool)).catch(fail);
    } catch (err) {
      fail(err);
    }
  }
}
