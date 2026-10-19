import { BatteryFull, Wifi } from "lucide-react"

import trayIcon from "@/assets/tray.png"
import { cn } from "@/lib/utils"

type MenuBarProps = {
  compact?: boolean
  className?: string
}

export function MenuBar({ compact = false, className }: MenuBarProps) {
  return (
    <div
      aria-hidden="true"
      className={cn(
        "flex items-center justify-between rounded-t-xl border border-white/10 border-b-white/[0.06] bg-[#182136]/90 px-3 shadow-[inset_0_1px_0_rgba(255,255,255,0.05)] backdrop-blur-xl",
        compact ? "h-8" : "h-9",
        className
      )}
    >
      <img
        src={trayIcon}
        alt=""
        width={36}
        height={36}
        className="size-3.5 opacity-90 invert"
      />
      <div className="flex items-center gap-2 text-white/55">
        <Wifi className="size-3" />
        <BatteryFull className="size-3.5" />
        <span
          className={cn(
            "font-medium text-white/70",
            compact ? "text-[10px]" : "text-[11px]"
          )}
        >
          9:41
        </span>
      </div>
    </div>
  )
}
