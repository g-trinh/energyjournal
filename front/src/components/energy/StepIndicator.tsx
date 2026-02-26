interface StepIndicatorProps {
  currentStep: 1 | 2
  onStepClick: (step: 1 | 2) => void
}

export default function StepIndicator({
  currentStep,
  onStepClick,
}: StepIndicatorProps) {
  const isStepOneActive = currentStep === 1
  const isStepTwoActive = currentStep === 2

  return (
    <nav className="energy-step-indicator" aria-label="Form progress">
      <div className="energy-step-indicator-inner">
        <div className="energy-step-indicator-line" aria-hidden="true" />

        <div className="energy-step-item energy-step-item-left">
          <button
            type="button"
            className={[
              'energy-step-dot',
              isStepOneActive ? 'is-active' : '',
              isStepTwoActive ? 'is-complete' : '',
            ]
              .filter(Boolean)
              .join(' ')}
            aria-label="Go to step 1"
            aria-current={isStepOneActive ? 'step' : undefined}
            onClick={() => onStepClick(1)}
          >
            {isStepTwoActive ? (
              <svg
                width="12"
                height="12"
                viewBox="0 0 12 12"
                aria-hidden="true"
                className="energy-step-check"
              >
                <path
                  d="M2 6.2 4.6 8.8 10 3.4"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="1.8"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            ) : (
              '1'
            )}
          </button>
          <span className="energy-step-label">ENERGY LEVELS</span>
        </div>

        <div className="energy-step-item energy-step-item-right">
          <button
            type="button"
            className={[
              'energy-step-dot',
              isStepTwoActive ? 'is-active' : '',
            ]
              .filter(Boolean)
              .join(' ')}
            aria-label="Go to step 2"
            aria-current={isStepTwoActive ? 'step' : undefined}
            onClick={() => onStepClick(2)}
            disabled={!isStepTwoActive}
          >
            2
          </button>
          <span
            className={[
              'energy-step-label',
              isStepTwoActive ? 'is-active' : '',
            ]
              .filter(Boolean)
              .join(' ')}
          >
            CONTEXT
          </span>
        </div>
      </div>
    </nav>
  )
}
