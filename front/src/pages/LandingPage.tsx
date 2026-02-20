import { Link } from 'react-router-dom'
import { Line, LineChart, ResponsiveContainer } from 'recharts'
import '../App.css'
import '@/styles/landing.css'

const DEMO_DATA = [
  { day: 'Mon', physical: 7, mental: 6, emotional: 8 },
  { day: 'Tue', physical: 8, mental: 5, emotional: 7 },
  { day: 'Wed', physical: 6, mental: 7, emotional: 8 },
  { day: 'Thu', physical: 7, mental: 6, emotional: 9 },
  { day: 'Fri', physical: 5, mental: 4, emotional: 6 },
  { day: 'Sat', physical: 8, mental: 7, emotional: 8 },
  { day: 'Sun', physical: 7, mental: 6, emotional: 8 },
]

export default function LandingPage() {
  return (
    <div className="app landing-page">
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <main aria-label="Landing page" className="landing-main">
        <section className="landing-hero" aria-labelledby="hero-heading">
          <div className="landing-hero-inner">
            <p className="landing-eyebrow">PERSONAL ENERGY INTELLIGENCE</p>
            <h1 id="hero-heading" className="landing-hero-title">
              Know your energy,
              {' '}
              <em className="landing-hero-title-accent">design your best days.</em>
            </h1>
            <p className="landing-hero-subline">
              Track Physical, Mental, and Emotional energy - discover the patterns that shape your
              performance, day after day.
            </p>
            <Link to="/timespending" className="landing-cta-primary">
              Start tracking →
            </Link>
          </div>
        </section>

        <div
          role="img"
          aria-label="7-day energy preview: Physical, Mental and Emotional trends"
          className="landing-chart-card landing-card"
        >
          <div className="landing-chart-legend" aria-hidden="true">
            <span className="landing-chart-legend-item">
              <span className="landing-chart-dot landing-chart-dot-physical" />
              Physical
            </span>
            <span className="landing-chart-legend-item">
              <span className="landing-chart-dot landing-chart-dot-mental" />
              Mental
            </span>
            <span className="landing-chart-legend-item">
              <span className="landing-chart-dot landing-chart-dot-emotional" />
              Emotional
            </span>
          </div>

          <div className="landing-chart-body">
            <ResponsiveContainer width="100%" height={160}>
              <LineChart data={DEMO_DATA}>
                <Line type="monotone" dataKey="physical" stroke="#c4826d" dot={false} strokeWidth={2} />
                <Line type="monotone" dataKey="mental" stroke="#7eb8b3" dot={false} strokeWidth={2} />
                <Line type="monotone" dataKey="emotional" stroke="#8fa58b" dot={false} strokeWidth={2} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

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
