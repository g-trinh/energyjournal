interface LogoMarkProps {
  size?: number
}

export default function LogoMark({ size = 36 }: LogoMarkProps) {
  return (
    <svg
      aria-hidden="true"
      className="topbar-logo-mark"
      viewBox="0 0 32 32"
      width={size}
      height={size}
    >
      <circle cx="16" cy="16" r="12" fill="none" stroke="currentColor" strokeWidth="1.5" opacity="0.3" />
      <circle cx="16" cy="16" r="8" fill="none" stroke="currentColor" strokeWidth="1.5" opacity="0.5" />
      <circle cx="16" cy="16" r="4" fill="currentColor" />
    </svg>
  )
}
