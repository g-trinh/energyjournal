import '../styles/auth.css'

export default function AuthPage() {
  return (
    <div className="app">
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <main className="auth-content">
        <header className="auth-header">
          <h1 className="auth-headline">Access your account</h1>
          <p className="auth-subline">
            Access your Energy Journal or create a new account in one place.
          </p>
        </header>

        <div className="auth-cards">
          <section className="auth-card">
            <h2 className="auth-card-title">Log In</h2>
            <p className="auth-card-description">Return to your dashboard</p>
            {/* LoginCard will be placed here in AUTH-003 */}
          </section>

          <section className="auth-card">
            <h2 className="auth-card-title">Sign Up</h2>
            <p className="auth-card-description">
              Create your account in seconds
            </p>
            {/* SignupCard will be placed here in AUTH-004 */}
          </section>
        </div>
      </main>
    </div>
  )
}
