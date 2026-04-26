<template>
  <main class="dashboard">
    <header class="dash-header">
      <div>
        <h1>Dashboard</h1>
        <p class="subtitle">Today is Yumsday!</p>
      </div>
    </header>

    <section>
      <div class="features-grid">
        <div v-for="feature in props.features" :key="feature.label" class="feature-card card">
          <span class="feature-icon">{{ feature.icon }}</span>
          <span class="feature-label">{{ feature.label }}</span>
        </div>
      </div>
    </section>
  </main>
</template>

<script setup>
const props = defineProps({
  features: {
    type: Array,
    default: () => [],
  },
})
</script>

<style scoped>
/* ── Layout ─────────────────────────────────────────────────── */
.dashboard {
  min-height: 100vh;
  padding: 2rem;
}

/* ── Header ───────────────────────────────────────────────────── */
.dash-header {
  max-width: 960px;
  margin: 0 auto 3rem;
}

.dash-header h1 {
  font-size: 2.5rem;
  margin: 0 0 0.5rem;
  color: #6A994E;
}

.subtitle {
  font-size: 0.9rem;
  opacity: 0.6;
  margin: 0;
}
</style>
