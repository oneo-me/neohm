import { ArrowUpRight, Download, ShieldCheck } from "lucide-react"

import appIcon from "@/assets/app-icon.webp"
import { Reveal } from "@/components/site/reveal"
import { SectionHeading } from "@/components/site/section-heading"
import { buttonVariants } from "@/components/ui/button"
import { useLanguage } from "@/lib/language"
import { DOWNLOAD_URL } from "@/lib/links"
import { cn } from "@/lib/utils"

export function Setup() {
  const { t } = useLanguage()

  return (
    <section
      id="setup"
      className="relative scroll-mt-24 border-t border-white/[0.05] py-24 sm:py-28"
    >
      <div className="mx-auto max-w-6xl px-5 sm:px-6">
        <Reveal>
          <SectionHeading
            eyebrow={t.setup.eyebrow}
            title={t.setup.title}
            body={t.setup.body}
          />
          <ul className="mt-6 flex flex-wrap gap-2">
            {t.setup.requirements.map((requirement) => (
              <li
                key={requirement}
                className="rounded-full border border-white/[0.08] bg-white/[0.03] px-3 py-1 font-mono text-[11px] text-muted-foreground"
              >
                {requirement}
              </li>
            ))}
          </ul>
        </Reveal>

        <div className="mt-12 grid items-start gap-6 lg:grid-cols-[0.95fr_1.05fr] lg:gap-10">
          <div className="flex flex-col gap-3">
            {t.setup.steps.map((step, index) => (
              <Reveal key={step.title} delay={index * 70}>
                <div className="flex gap-4 rounded-2xl border border-white/[0.06] bg-card/40 p-5 transition-colors hover:border-white/[0.1]">
                  <span className="flex size-8 shrink-0 items-center justify-center rounded-full bg-white/[0.05] font-mono text-xs text-brand-blue ring-1 ring-white/10">
                    {index + 1}
                  </span>
                  <div>
                    <h3 className="text-[15px] font-medium tracking-tight">
                      {step.title}
                    </h3>
                    <p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">
                      {step.body}
                    </p>
                  </div>
                </div>
              </Reveal>
            ))}
          </div>

          <Reveal delay={120} className="flex flex-col gap-4">
            <div className="relative overflow-hidden rounded-2xl border border-white/[0.09] bg-[#050a13] shadow-[0_24px_60px_rgba(1,5,14,0.45)]">
              <div
                aria-hidden="true"
                className="pointer-events-none absolute -top-28 right-[-20%] size-72 rounded-full bg-[radial-gradient(circle,rgba(81,133,255,0.26),transparent_66%)] blur-2xl"
              />
              <div
                aria-hidden="true"
                className="pointer-events-none absolute -bottom-36 left-[-18%] size-72 rounded-full bg-[radial-gradient(circle,rgba(246,195,126,0.1),transparent_66%)] blur-2xl"
              />

              <div className="relative p-6 sm:p-7">
                <div className="flex items-center gap-4">
                  <img
                    src={appIcon}
                    alt=""
                    width={512}
                    height={512}
                    loading="lazy"
                    className="size-14 rounded-2xl shadow-[0_12px_30px_rgba(0,0,0,0.45)]"
                  />
                  <div>
                    <p className="text-[10px] font-medium tracking-[0.16em] text-brand-mint uppercase">
                      {t.setup.free}
                    </p>
                    <h3 className="mt-1 text-lg font-medium tracking-tight">
                      {t.setup.downloadTitle}
                    </h3>
                    <p className="mt-1 text-sm text-muted-foreground">
                      {t.setup.downloadBody}
                    </p>
                  </div>
                </div>

                <a
                  href={DOWNLOAD_URL}
                  target="_blank"
                  rel="noreferrer"
                  className={cn(
                    buttonVariants({ size: "lg" }),
                    "mt-6 h-11 w-full rounded-xl px-5 text-sm font-semibold"
                  )}
                >
                  <Download data-icon="inline-start" />
                  {t.setup.download}
                </a>

                <a
                  href={DOWNLOAD_URL}
                  target="_blank"
                  rel="noreferrer"
                  className="mt-3 inline-flex items-center gap-1 text-xs text-muted-foreground transition-colors hover:text-foreground"
                >
                  {t.setup.releases}
                  <ArrowUpRight className="size-3.5" />
                </a>

                <div className="mt-5 border-t border-white/[0.06] pt-4">
                  <ul className="flex flex-wrap gap-1.5">
                    {t.setup.platforms.map((platform) => (
                      <li
                        key={platform.name}
                        className="inline-flex items-center gap-1.5 rounded-md border border-white/[0.08] bg-white/[0.03] px-2 py-1 text-[11px]"
                      >
                        <span className="text-foreground/80">
                          {platform.name}
                        </span>
                        <span className="text-muted-foreground">
                          {platform.status}
                        </span>
                      </li>
                    ))}
                  </ul>
                  <p className="mt-2.5 text-[11px] text-muted-foreground">
                    {t.setup.architectures}
                  </p>
                </div>
              </div>
            </div>

            <div className="flex gap-3.5 rounded-2xl border border-brand-mint/15 bg-brand-mint/[0.04] p-5">
              <ShieldCheck className="mt-0.5 size-5 shrink-0 text-brand-mint" />
              <div>
                <h3 className="text-sm font-medium">{t.setup.privacyTitle}</h3>
                <p className="mt-1.5 text-[13px] leading-relaxed text-muted-foreground">
                  {t.setup.privacyBody}
                </p>
              </div>
            </div>
          </Reveal>
        </div>
      </div>
    </section>
  )
}
