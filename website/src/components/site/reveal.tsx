import { useEffect, useRef, useState, type ReactNode } from "react"

import { cn } from "@/lib/utils"

type RevealProps = {
  children: ReactNode
  className?: string
  delay?: number
}

export function Reveal({ children, className, delay = 0 }: RevealProps) {
  const ref = useRef<HTMLDivElement>(null)
  const [shown, setShown] = useState(() =>
    typeof window === "undefined"
      ? true
      : window.matchMedia("(prefers-reduced-motion: reduce)").matches
  )

  useEffect(() => {
    const element = ref.current
    if (
      !element ||
      window.matchMedia("(prefers-reduced-motion: reduce)").matches
    ) {
      return
    }

    const reveal = () => setShown(true)
    const bounds = element.getBoundingClientRect()

    if (bounds.top < window.innerHeight * 0.95 && bounds.bottom > 0) {
      const frame = window.setTimeout(reveal, 60)
      return () => window.clearTimeout(frame)
    }

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) {
          reveal()
          observer.disconnect()
        }
      },
      { rootMargin: "0px 0px -8% 0px", threshold: 0.05 }
    )

    observer.observe(element)
    const fallback = window.setTimeout(reveal, 1500)

    return () => {
      observer.disconnect()
      window.clearTimeout(fallback)
    }
  }, [])

  return (
    <div
      ref={ref}
      className={cn(
        "transition-all duration-700 ease-out",
        shown ? "translate-y-0 opacity-100" : "translate-y-5 opacity-0",
        className
      )}
      style={{ transitionDelay: `${delay}ms` }}
    >
      {children}
    </div>
  )
}
