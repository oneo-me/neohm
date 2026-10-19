import { ArrowRight } from "lucide-react"

import appIcon from "@/assets/app-icon.webp"
import { GithubIcon } from "@/components/site/github-icon"
import { Reveal } from "@/components/site/reveal"
import { buttonVariants } from "@/components/ui/button"
import { useLanguage } from "@/lib/language"
import { DOWNLOAD_URL, REPO_URL } from "@/lib/links"
import { cn } from "@/lib/utils"

export function CallToAction() {
  const { t } = useLanguage()

  return (
    <section className="pb-24 sm:pb-28">
      <div className="mx-auto max-w-6xl px-5 sm:px-6">
        <Reveal>
          <div className="relative overflow-hidden rounded-3xl border border-white/[0.09] bg-[#091120] px-6 py-16 text-center sm:px-12 sm:py-20">
            <div
              aria-hidden="true"
              className="pointer-events-none absolute inset-0"
            >
              <div className="absolute inset-0 bg-grid-faint [mask-image:radial-gradient(ellipse_65%_80%_at_50%_50%,black,transparent)] opacity-60" />
              <div className="absolute -top-40 left-1/2 h-96 w-[680px] -translate-x-1/2 rounded-full bg-[radial-gradient(ellipse_at_center,rgba(81,133,255,0.3),transparent_65%)] blur-2xl" />
              <div className="absolute -right-20 -bottom-32 size-80 rounded-full bg-[radial-gradient(circle,rgba(246,195,126,0.14),transparent_65%)] blur-2xl" />
            </div>

            <div className="relative flex flex-col items-center">
              <img
                src={appIcon}
                alt=""
                width={512}
                height={512}
                loading="lazy"
                className="size-14 rounded-2xl shadow-[0_16px_40px_rgba(0,0,0,0.45)]"
              />
              <h2 className="mt-7 max-w-2xl text-3xl font-semibold tracking-[-0.03em] text-balance sm:text-4xl lg:text-[2.8rem] lg:leading-[1.08]">
                {t.cta.title.map((line) => (
                  <span key={line} className="block">
                    {line}
                  </span>
                ))}
              </h2>
              <p className="mt-4 max-w-xl text-[15px] leading-relaxed text-muted-foreground">
                {t.cta.body}
              </p>
              <div className="mt-8 flex flex-col gap-3 sm:flex-row">
                <a
                  href={DOWNLOAD_URL}
                  target="_blank"
                  rel="noreferrer"
                  className={cn(
                    buttonVariants({ size: "lg" }),
                    "h-11 rounded-xl px-5 text-sm font-semibold"
                  )}
                >
                  {t.cta.primary}
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
                  {t.cta.secondary}
                </a>
              </div>
            </div>
          </div>
        </Reveal>
      </div>
    </section>
  )
}
