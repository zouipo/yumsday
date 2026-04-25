<template>
  <main class="home">
    <!-- ─── GUEST VIEW ─── -->
    <section v-if="!auth.isAuthenticated" class="guest">
      <div class="hero">
        <span class="eyebrow">Welcome</span>
        <h1>Today is<br /><em>Yumsday!</em></h1>
      </div>

      <div class="features">
        <div v-for="f in features" :key="f.title" class="feature-card card">
          <span class="feature-icon">{{ f.icon }}</span>
          <h3>{{ f.title }}</h3>
          <p>{{ f.description }}</p>
        </div>
      </div>
    </section>

    <!-- ─── AUTHENTICATED VIEW ─── -->
    <section v-else class="authenticated">
      <DashboardPage :features="features" />
    </section>
  </main>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import DashboardPage from './DashboardPage.vue'

const auth = useAuthStore()

// ── Shared data ─────────────────────────────────────────────
const features = [
  { icon: '📆', title: 'Plan out your meals',  label: 'Menu',         description: 'An interactive and handy menu.' },
  { icon: '🍲', title: 'Save your recipes',    label: 'Recipes',      description: 'A true personal cookbook to always remind you of your grandmother\'s delicious recipes.' },
  { icon: '📋', title: 'Go grocery shopping',  label: 'Grocery list', description: 'Simplify your grocery shopping and make it less of a hassle.' },
]
</script>

<style scoped>
/* ── Layout ─────────────────────────────────────────────────── */
.home { min-height: 100vh; padding: 2rem; }

/* ── Guest ───────────────────────────────────────────────────── */
.hero { max-width: 680px; margin: 6rem auto 4rem; text-align: center; }
.eyebrow { font-size: .75rem; letter-spacing: .2em; text-transform: uppercase; opacity: .5; }
.hero h1 { font-size: clamp(2.5rem, 6vw, 4rem); line-height: 1.1; margin: .5rem 0 1.25rem; color: var(--color-primary); }
.hero h1 em { font-style: italic; color: var(--color-danger); }
.cta-group { display: flex; gap: 1rem; justify-content: center; flex-wrap: wrap; }

.features { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.5rem; max-width: 860px; margin: 0 auto; }
.feature-card { padding: 1.5rem; }
.feature-icon { font-size: 1.75rem; }
.feature-card h3 { margin: .5rem 0 .25rem; }
.feature-card p  { opacity: .6; margin: 0; }

/* ── Authenticated ──────────────────────────────────────────── */
.dash-header { display: flex; justify-content: space-between; align-items: flex-end;
  max-width: 960px; margin: 0 auto 2.5rem; }
.dash-header h1 { margin: 0; font-size: 2rem; }

.features-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 1rem; max-width: 960px; margin: 0 auto 2.5rem; }
.feature-card { padding: 1.25rem 1.5rem; display: flex; flex-direction: column; gap: .25rem; }
.feature-label { font-size: .8rem; opacity: .5; text-transform: uppercase; letter-spacing: .08em; }
</style>