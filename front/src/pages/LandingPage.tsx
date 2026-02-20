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
          <p className="landing-section-sub">Rate each one from 0 to 10, every day.</p>
          <div className="landing-features-grid">
            <article className="landing-card landing-feature-card landing-feature-card-physical">
              <svg aria-hidden="true" className="landing-feature-icon" viewBox="0 0 24 24">
                <path d="M12 4v16M4 12h16" fill="none" stroke="currentColor" strokeWidth="1.7" strokeLinecap="round" />
              </svg>
              <h3 className="landing-feature-title">Physical</h3>
              <p className="landing-feature-tag">BODY · ENDURANCE</p>
              <p className="landing-feature-desc">
                Follow how sleep, activity, and recovery shape your day-to-day stamina.
              </p>
            </article>
            <article className="landing-card landing-feature-card landing-feature-card-mental">
              <svg aria-hidden="true" className="landing-feature-icon" viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="7.5" fill="none" stroke="currentColor" strokeWidth="1.7" />
                <circle cx="12" cy="12" r="3.2" fill="none" stroke="currentColor" strokeWidth="1.7" />
              </svg>
              <h3 className="landing-feature-title">Mental</h3>
              <p className="landing-feature-tag">FOCUS · CLARITY</p>
              <p className="landing-feature-desc">
                Understand when concentration peaks and when your cognitive load becomes too high.
              </p>
            </article>
            <article className="landing-card landing-feature-card landing-feature-card-emotional">
              <svg aria-hidden="true" className="landing-feature-icon" viewBox="0 0 24 24">
                <path d="M12 20s-7-4.4-7-10a4 4 0 0 1 7-2.6A4 4 0 0 1 19 10c0 5.6-7 10-7 10Z" fill="none" stroke="currentColor" strokeWidth="1.7" />
              </svg>
              <h3 className="landing-feature-title">Emotional</h3>
              <p className="landing-feature-tag">MOOD · RESILIENCE</p>
              <p className="landing-feature-desc">
                Spot mood trends early and build routines that protect calm, confidence, and resilience.
              </p>
            </article>
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
