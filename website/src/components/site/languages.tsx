import { Globe2 } from "lucide-react"

import { Reveal } from "@/components/site/reveal"
import { languageOptions } from "@/lib/copy"
import { useLanguage } from "@/lib/language"

export function Languages() {
  const { t } = useLanguage()

  return (
    <section aria-labelledby="languages-title" className="py-16 sm:py-20">
      <div className="mx-auto max-w-6xl px-5 sm:px-6">
        <Reveal>
          <div className="flex flex-col gap-8 rounded-3xl border border-white/[0.07] bg-card/50 px-6 py-8 sm:px-8 lg:flex-row lg:items-center lg:justify-between">
            <div className="flex items-start gap-4">
              <div className="flex size-11 shrink-0 items-center justify-center rounded-2xl border border-brand-mint/20 bg-brand-mint/10 text-brand-mint">
                <Globe2 className="size-5" />
              </div>
              <div>
                <p className="text-[10px] font-medium tracking-[0.16em] text-muted-foreground uppercase">
                  {t.languages.eyebrow}
                </p>
                <h2
                  id="languages-title"
                  className="mt-1.5 text-xl font-semibold tracking-tight sm:text-2xl"
                >
                  {t.languages.title}
                </h2>
                <p className="mt-2 max-w-md text-sm leading-relaxed text-muted-foreground">
                  {t.languages.body}
                </p>
              </div>
            </div>

            <ul className="flex flex-wrap gap-2 lg:justify-end">
              {languageOptions.map((language) => (
                <li
                  key={language.value}
                  className="rounded-full border border-white/[0.08] bg-white/[0.03] px-3.5 py-1.5 text-sm text-foreground/80 transition-colors hover:border-white/15 hover:text-foreground"
                >
                  {language.label}
                </li>
              ))}
            </ul>
          </div>
        </Reveal>
      </div>
    </section>
  )
}
