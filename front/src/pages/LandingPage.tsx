import '../App.css'
import '@/styles/landing.css'

export default function LandingPage() {
  return (
    <div className="app landing-page">
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <main aria-label="Landing page" className="landing-main">
        <section className="landing-hero" aria-labelledby="hero-heading">
          <div className="landing-hero-inner">
            <h1 id="hero-heading">Know your energy</h1>
          </div>
        </section>

        <section className="landing-features" aria-labelledby="features-heading">
          <h2 id="features-heading">Three dimensions of energy</h2>
          <div className="landing-features-grid">
            <article className="landing-card">Physical</article>
            <article className="landing-card">Mental</article>
            <article className="landing-card">Emotional</article>
          </div>
        </section>

        <section className="landing-steps" aria-labelledby="steps-heading">
          <h2 id="steps-heading">Simple to start, powerful over time</h2>
        </section>

        <section className="landing-cta-section" aria-labelledby="cta-heading">
          <h2 id="cta-heading">Start understanding your energy today</h2>
          <div className="landing-cta-banner" />
        </section>
      </main>

      <footer className="landing-footer">© 2026 Energy Journal</footer>
    </div>
  )
}
