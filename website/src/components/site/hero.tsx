import { ArrowRight, Sparkles } from "lucide-react"

import { GithubIcon } from "@/components/site/github-icon"
import { MenuBar } from "@/components/site/menu-bar"
import { Reveal } from "@/components/site/reveal"
import { Badge } from "@/components/ui/badge"
import { buttonVariants } from "@/components/ui/button"
import { useLanguage } from "@/lib/language"
import { DOWNLOAD_URL, REPO_URL } from "@/lib/links"
import { cn } from "@/lib/utils"

export function Hero() {
  const { t, screens, lang } = useLanguage()

  return (
    <section className="relative overflow-hidden pt-12 pb-20 sm:pt-16 lg:pt-24 lg:pb-32">
      <div className="mx-auto grid max-w-6xl items-center gap-16 px-5 sm:px-6 lg:grid-cols-[1.04fr_0.96fr] lg:gap-8">
        <div className="relative z-10 max-w-xl">
          <Reveal>
            <Badge
              variant="outline"
              className="h-auto gap-2 rounded-full border-white/10 bg-white/[0.03] px-3 py-1.5 text-xs font-normal text-muted-foreground backdrop-blur"
            >
              <Sparkles className="text-brand-amber" />
              {t.hero.badge}
            </Badge>
          </Reveal>

          <Reveal delay={70}>
            <h1 className="mt-6 text-[2.75rem] leading-[0.98] font-semibold tracking-[-0.045em] text-balance sm:text-6xl lg:text-[4.15rem]">
              {t.hero.titleLead}
              <br />
              <span className="text-gradient">{t.hero.titleAccent}</span>
            </h1>
          </Reveal>

          <Reveal delay={140}>
            <p
              className={cn(
                "mt-6 text-[15px] leading-relaxed text-muted-foreground sm:text-lg sm:leading-relaxed",
                lang === "en" ? "max-w-lg" : "max-w-[34rem] text-pretty"
              )}
            >
              {t.hero.body}
            </p>
          </Reveal>

          <Reveal delay={210}>
            <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:items-center">
              <a
                href={DOWNLOAD_URL}
                target="_blank"
                rel="noreferrer"
                className={cn(
                  buttonVariants({ size: "lg" }),
                  "h-11 rounded-xl px-5 text-sm font-semibold"
                )}
              >
                {t.hero.primary}
                <ArrowRight data-icon="inline-end" />
              </a>
              <a
                href={REPO_URL}
                target="_blank"
                rel="noreferrer"
                className={cn(
                  buttonVariants({ variant: "outline", size: "lg" }),
                  "h-11 rounded-xl px-5 text-sm"
                )}
              >
                <GithubIcon className="size-4" />
                {t.hero.secondary}
              </a>
            </div>
          </Reveal>

          <Reveal delay={280}>
            <div className="mt-8 flex flex-wrap items-center gap-x-5 gap-y-2 text-xs text-muted-foreground">
              {t.hero.meta.map((item) => (
                <span key={item}>{item}</span>
              ))}
            </div>
            <p className="mt-4 flex items-center gap-2 text-xs text-muted-foreground/90">
              <span className="size-1.5 rounded-full bg-brand-amber shadow-[0_0_8px_2px_rgba(246,195,126,0.35)]" />
              {t.hero.note}
            </p>
          </Reveal>
        </div>

        <div className="relative mx-auto w-full max-w-[404px]">
          <div
            aria-hidden="true"
            className="pointer-events-none absolute -inset-x-24 -top-20 bottom-0 -z-10"
          >
            <div className="absolute top-16 left-1/2 size-[430px] -translate-x-1/2 animate-[glow-drift_11s_ease-in-out_infinite] rounded-full bg-[radial-gradient(circle,rgba(70,124,255,0.28),transparent_66%)] blur-2xl" />
            <div className="absolute top-1/2 -right-8 size-64 rounded-full bg-[radial-gradient(circle,rgba(246,195,126,0.13),transparent_66%)] blur-2xl" />
            <div className="absolute -bottom-10 -left-10 size-60 rounded-full bg-[radial-gradient(circle,rgba(183,156,255,0.12),transparent_66%)] blur-2xl" />
          </div>

          <Reveal delay={160} className="relative">
            <div className="relative">
              <MenuBar className="relative z-10" />
              <img
                src={screens.home}
                alt={t.alts.home}
                width={776}
                height={1144}
                className="relative z-0 block w-full drop-shadow-[0_36px_70px_rgba(1,5,14,0.6)] select-none"
                fetchPriority="high"
                draggable={false}
              />
            </div>
          </Reveal>
        </div>
      </div>
    </section>
  )
}
