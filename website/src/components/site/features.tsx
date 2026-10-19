import {
  Blinds,
  Blocks,
  ChevronRight,
  Fan,
  Globe2,
  Lamp,
  LayoutGrid,
  Lightbulb,
  Lock,
  Moon,
  MousePointer2,
  MousePointerClick,
  PanelTop,
  Plug,
  Settings2,
  Snowflake,
  Tv,
  Wind,
  type LucideIcon,
} from "lucide-react"

import { MenuBar } from "@/components/site/menu-bar"
import { Reveal } from "@/components/site/reveal"
import { SectionHeading } from "@/components/site/section-heading"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useLanguage } from "@/lib/language"
import { cn } from "@/lib/utils"

const icons: LucideIcon[] = [
  LayoutGrid,
  MousePointerClick,
  Snowflake,
  PanelTop,
  Settings2,
  Blocks,
]

const iconStyles = [
  "border-brand-blue/20 bg-brand-blue/10 text-brand-blue",
  "border-brand-violet/20 bg-brand-violet/10 text-brand-violet",
  "border-brand-amber/20 bg-brand-amber/10 text-brand-amber",
  "border-brand-amber/20 bg-brand-amber/10 text-brand-amber",
  "border-brand-mint/20 bg-brand-mint/10 text-brand-mint",
  "border-brand-pink/20 bg-brand-pink/10 text-brand-pink",
]

function MiniTiles() {
  const tiles = [
    {
      icon: Lightbulb,
      tone: "text-brand-amber",
      lit: "border-brand-amber/25 bg-[radial-gradient(120%_120%_at_20%_0%,rgba(246,195,126,0.18),rgba(11,19,34,0.35))]",
    },
    {
      icon: Snowflake,
      tone: "text-brand-blue",
      lit: "border-brand-blue/25 bg-[radial-gradient(120%_120%_at_20%_0%,rgba(125,178,255,0.18),rgba(11,19,34,0.35))]",
    },
    {
      icon: Blinds,
      tone: "text-brand-violet",
      lit: "border-brand-violet/25 bg-[radial-gradient(120%_120%_at_20%_0%,rgba(183,156,255,0.18),rgba(11,19,34,0.35))]",
    },
    {
      icon: Lamp,
      tone: "text-white/45",
      lit: "border-white/[0.06] bg-white/[0.02]",
    },
  ]

  return (
    <div aria-hidden="true" className="grid grid-cols-2 gap-2">
      {tiles.map((tile, index) => (
        <div
          key={index}
          className={cn(
            "rounded-xl border p-2.5 transition-colors",
            tile.lit,
            index === 0 &&
              "motion-safe:animate-[tile-breathe_5s_ease-in-out_infinite]"
          )}
        >
          <div className="flex items-center justify-between">
            <span
              className={cn(
                "flex size-6 items-center justify-center rounded-full bg-black/20 ring-1 ring-white/5",
                tile.tone
              )}
            >
              <tile.icon className="size-3" />
            </span>
            <ChevronRight className="size-3 text-white/20" />
          </div>
          <div className="mt-3 h-1.5 w-3/5 rounded-full bg-white/15" />
          <div className="mt-1 h-1 w-2/5 rounded-full bg-white/10" />
        </div>
      ))}
    </div>
  )
}

function MiniSwitch() {
  return (
    <div aria-hidden="true" className="grid grid-cols-[1.15fr_1fr] gap-2">
      <div className="relative rounded-xl border border-brand-amber/20 bg-[radial-gradient(120%_120%_at_20%_0%,rgba(246,195,126,0.16),rgba(11,19,34,0.35))] p-3">
        <div className="flex items-center justify-between">
          <span className="flex size-8 items-center justify-center rounded-full bg-brand-amber/15 text-brand-amber ring-1 ring-brand-amber/20">
            <Lightbulb className="size-4" />
          </span>
          <span className="flex size-5 items-center justify-center rounded-full bg-white/[0.05] ring-1 ring-white/10">
            <ChevronRight className="size-3 text-white/40" />
          </span>
        </div>
        <div className="mt-3 h-1.5 w-2/5 rounded-full bg-white/20" />
        <div className="mt-1.5 h-1 w-1/4 rounded-full bg-white/10" />
        <MousePointer2 className="absolute right-3 bottom-3 size-4 fill-white text-[#0b1322] drop-shadow-[0_2px_4px_rgba(0,0,0,0.55)]" />
      </div>
      <div className="rounded-xl border border-white/[0.06] bg-white/[0.02] p-3">
        <div className="flex items-center justify-between">
          <span className="flex size-8 items-center justify-center rounded-full bg-white/[0.04] text-white/35 ring-1 ring-white/[0.06]">
            <Plug className="size-4" />
          </span>
          <ChevronRight className="size-3 text-white/20" />
        </div>
        <div className="mt-3 h-1.5 w-1/2 rounded-full bg-white/10" />
        <div className="mt-1.5 h-1 w-1/3 rounded-full bg-white/[0.07]" />
      </div>
    </div>
  )
}

function MiniSettings() {
  const pills = ["EN", "简", "繁", "日", "한"]

  return (
    <div
      aria-hidden="true"
      className="rounded-xl border border-white/[0.06] bg-black/20 p-3"
    >
      <div className="flex items-center justify-between gap-2">
        <div className="flex gap-1.5">
          {pills.map((pill, index) => (
            <span
              key={pill}
              className={cn(
                "rounded-full px-2 py-0.5 text-[9px] ring-1",
                index === 0
                  ? "bg-brand-mint/15 text-brand-mint ring-brand-mint/20"
                  : "bg-white/[0.03] text-muted-foreground ring-white/[0.06]"
              )}
            >
              {pill}
            </span>
          ))}
        </div>
        <Globe2 className="size-3.5 text-white/30" />
      </div>
      <div className="mt-2 h-5 rounded-md border border-white/[0.06] bg-white/[0.02]" />
    </div>
  )
}

function MiniDevices() {
  const devices = [
    { icon: Lightbulb, tone: "text-brand-amber" },
    { icon: Fan, tone: "text-brand-blue" },
    { icon: Lock, tone: "text-brand-violet" },
    { icon: Tv, tone: "text-brand-pink" },
  ]

  return (
    <div
      aria-hidden="true"
      className="flex items-center justify-between gap-2 rounded-xl border border-white/[0.06] bg-black/20 p-3"
    >
      {devices.map((device, index) => (
        <span
          key={index}
          className="flex flex-1 items-center justify-center rounded-lg bg-white/[0.04] py-2.5 ring-1 ring-white/[0.06]"
        >
          <device.icon className={cn("size-4", device.tone)} />
        </span>
      ))}
    </div>
  )
}

function MiniClimate() {
  const chips = [
    { icon: Snowflake, tone: "text-brand-blue" },
    { icon: Fan, tone: "text-brand-blue" },
    { icon: Moon, tone: "text-muted-foreground" },
    { icon: Wind, tone: "text-brand-blue" },
  ]

  return (
    <div
      aria-hidden="true"
      className="flex items-center justify-between gap-3 rounded-xl border border-white/[0.06] bg-black/20 px-3 py-2.5"
    >
      <span className="text-2xl font-light tracking-tight text-white/90">
        16
        <span className="ml-0.5 align-top text-xs text-muted-foreground">
          °C
        </span>
      </span>
      <div className="flex gap-1.5">
        {chips.map((chip, index) => (
          <span
            key={index}
            className={cn(
              "flex size-8 items-center justify-center rounded-lg bg-white/[0.04] ring-1 ring-white/[0.06]",
              chip.tone
            )}
          >
            <chip.icon className="size-4" />
          </span>
        ))}
      </div>
      <div className="mt-2.5 flex items-center gap-2">
        <span className="text-[10px] text-white/30">−</span>
        <div className="relative h-1 flex-1 rounded-full bg-white/10">
          <div className="h-1 w-2/5 rounded-full bg-brand-blue/80" />
          <div className="absolute top-1/2 left-[38%] size-3 -translate-y-1/2 rounded-full bg-white shadow-[0_1px_4px_rgba(0,0,0,0.4)]" />
        </div>
        <span className="text-[10px] text-white/30">+</span>
      </div>
    </div>
  )
}

function visualFor(index: number) {
  if (index === 1) {
    return <MiniSwitch />
  }

  if (index === 2) {
    return <MiniClimate />
  }

  if (index === 3) {
    return (
      <div
        aria-hidden="true"
        className="rounded-xl bg-black/20 p-2.5 ring-1 ring-white/[0.04]"
      >
        <MenuBar compact className="rounded-lg" />
      </div>
    )
  }

  if (index === 4) {
    return <MiniSettings />
  }

  if (index === 5) {
    return <MiniDevices />
  }

  return <MiniTiles />
}

export function Features() {
  const { t } = useLanguage()

  return (
    <section id="features" className="relative scroll-mt-24 py-24 sm:py-28">
      <div className="mx-auto max-w-6xl px-5 sm:px-6">
        <Reveal>
          <SectionHeading
            eyebrow={t.features.eyebrow}
            title={t.features.title}
            body={t.features.body}
          />
        </Reveal>

        <div className="mt-12 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {t.features.items.map((item, index) => {
            const Icon = icons[index] ?? LayoutGrid
            const visual = visualFor(index)

            return (
              <Reveal key={item.title} delay={index * 55} className="h-full">
                <Card className="h-full transition-all duration-300 [--card-spacing:--spacing(5)] hover:-translate-y-0.5 hover:ring-white/[0.14]">
                  <CardHeader>
                    <div className="mb-2 flex items-start justify-between">
                      <div
                        className={cn(
                          "flex size-10 items-center justify-center rounded-xl border ring-1 ring-white/[0.04]",
                          iconStyles[index]
                        )}
                      >
                        <Icon className="size-[18px]" />
                      </div>
                      <span className="font-mono text-[10px] tracking-widest text-muted-foreground/40">
                        0{index + 1}
                      </span>
                    </div>
                    <CardTitle
                      role="heading"
                      aria-level={3}
                      className="text-[17px]"
                    >
                      {item.title}
                    </CardTitle>
                    <CardDescription className="text-[13.5px] leading-relaxed">
                      {item.body}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="mt-auto pt-2">{visual}</CardContent>
                </Card>
              </Reveal>
            )
          })}
        </div>
      </div>
    </section>
  )
}
