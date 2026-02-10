import * as React from "react"

const Input = React.forwardRef<
  HTMLInputElement,
  React.ComponentProps<"input">
>(({ className, type, ...props }, ref) => (
  <input type={type} className={className} ref={ref} {...props} />
))
Input.displayName = "Input"

export { Input }
