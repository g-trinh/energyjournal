import { useCallback, useEffect, useId, useMemo, useRef, useState } from 'react'
import type { KeyboardEvent } from 'react'

export interface ContextSelectOption {
  value: string
  label: string
}

interface ContextSelectProps {
  label: string
  options: ContextSelectOption[]
  value: string
  onChange: (value: string) => void
  disabled?: boolean
  isOpen?: boolean
  onOpenChange?: (open: boolean) => void
}

export default function ContextSelect({
  label,
  options,
  value,
  onChange,
  disabled = false,
  isOpen,
  onOpenChange,
}: ContextSelectProps) {
  const containerRef = useRef<HTMLDivElement | null>(null)
  const listboxId = useId()
  const isControlled = typeof isOpen === 'boolean'
  const [internalOpen, setInternalOpen] = useState(false)
  const [activeIndex, setActiveIndex] = useState(-1)

  const open = isControlled ? Boolean(isOpen) : internalOpen

  const selectedIndex = useMemo(
    () => options.findIndex((option) => option.value === value),
    [options, value],
  )

  const displayedLabel =
    options.find((option) => option.value === value)?.label ?? '\u2014'

  const setOpen = useCallback(
    (nextOpen: boolean) => {
      if (!isControlled) {
        setInternalOpen(nextOpen)
      }
      onOpenChange?.(nextOpen)
    },
    [isControlled, onOpenChange],
  )

  function openSelect() {
    if (disabled) {
      return
    }
    const nextIndex = selectedIndex >= 0 ? selectedIndex : 0
    setActiveIndex(nextIndex)
    setOpen(true)
  }

  const closeSelect = useCallback(() => {
    setOpen(false)
    setActiveIndex(-1)
  }, [setOpen])

  function moveActive(delta: number) {
    if (options.length === 0) {
      return
    }

    const baseIndex = activeIndex >= 0 ? activeIndex : selectedIndex >= 0 ? selectedIndex : 0
    const nextIndex = (baseIndex + delta + options.length) % options.length
    setActiveIndex(nextIndex)
  }

  function handleSelect(nextValue: string) {
    if (disabled) {
      return
    }
    onChange(nextValue)
    closeSelect()
  }

  function handleTriggerKeyDown(event: KeyboardEvent<HTMLButtonElement>) {
    if (disabled) {
      return
    }

    if (event.key === 'ArrowDown') {
      event.preventDefault()
      if (!open) {
        openSelect()
        return
      }
      moveActive(1)
      return
    }

    if (event.key === 'ArrowUp') {
      event.preventDefault()
      if (!open) {
        openSelect()
        return
      }
      moveActive(-1)
      return
    }

    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault()
      if (!open) {
        openSelect()
        return
      }
      if (activeIndex >= 0) {
        handleSelect(options[activeIndex].value)
      }
      return
    }

    if (event.key === 'Escape' && open) {
      event.preventDefault()
      closeSelect()
    }
  }

  function handleListboxKeyDown(event: KeyboardEvent<HTMLDivElement>) {
    if (!open) {
      return
    }

    if (event.key === 'ArrowDown') {
      event.preventDefault()
      moveActive(1)
      return
    }

    if (event.key === 'ArrowUp') {
      event.preventDefault()
      moveActive(-1)
      return
    }

    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault()
      if (activeIndex >= 0) {
        handleSelect(options[activeIndex].value)
      }
      return
    }

    if (event.key === 'Escape') {
      event.preventDefault()
      closeSelect()
    }
  }

  useEffect(() => {
    if (!open) {
      return
    }

    function handleClickOutside(event: MouseEvent) {
      if (!containerRef.current) {
        return
      }
      if (containerRef.current.contains(event.target as Node)) {
        return
      }
      closeSelect()
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [closeSelect, open])

  return (
    <div className="energy-context-select" ref={containerRef}>
      <label className="energy-context-label">{label}</label>
      <button
        type="button"
        className="energy-context-trigger"
        role="combobox"
        aria-expanded={open}
        aria-haspopup="listbox"
        aria-controls={listboxId}
        aria-label={label}
        disabled={disabled}
        onClick={() => {
          if (open) {
            closeSelect()
            return
          }
          openSelect()
        }}
        onKeyDown={handleTriggerKeyDown}
      >
        <span
          className={[
            'energy-context-value',
            value === '' ? 'is-placeholder' : '',
          ]
            .filter(Boolean)
            .join(' ')}
        >
          {displayedLabel}
        </span>
        <svg
          className={['energy-context-chevron', open ? 'is-open' : '']
            .filter(Boolean)
            .join(' ')}
          width="14"
          height="14"
          viewBox="0 0 20 20"
          aria-hidden="true"
        >
          <path
            d="M5 7l5 6 5-6"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.6"
            strokeLinecap="round"
          />
        </svg>
      </button>

      {open && (
        <div
          id={listboxId}
          className="energy-context-options"
          role="listbox"
          aria-label={label}
          onKeyDown={handleListboxKeyDown}
        >
          {options.map((option, index) => {
            const isSelected = option.value === value
            const isActive = index === activeIndex

            return (
              <button
                key={option.value || '__empty'}
                type="button"
                role="option"
                aria-selected={isSelected}
                className={[
                  'energy-context-option',
                  isSelected ? 'is-selected' : '',
                  isActive ? 'is-active' : '',
                ]
                  .filter(Boolean)
                  .join(' ')}
                onMouseEnter={() => setActiveIndex(index)}
                onMouseDown={(event) => {
                  event.preventDefault()
                }}
                onClick={() => handleSelect(option.value)}
              >
                {option.label}
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}
