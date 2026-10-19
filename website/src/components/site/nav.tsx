import { ArrowRight } from "lucide-react"

import { GithubIcon } from "@/components/site/github-icon"
import { LanguageSwitch } from "@/components/site/language-switch"
import { Logo } from "@/components/site/logo"
import { buttonVariants } from "@/components/ui/button"
import { useLanguage } from "@/lib/language"
import { DOWNLOAD_URL, REPO_URL } from "@/lib/links"
import { cn } from "@/lib/utils"

export function Nav() {
  const { t } = useLanguage()

  const links = [
    { href: "#features", label: t.nav.features },
    { href: "#screens", label: t.nav.screens },
    { href: "#setup", label: t.nav.setup },
    { href: "#faq", label: t.nav.faq },
  ]

  return (
    <header className="sticky top-0 z-50 border-b border-white/[0.06] bg-background/70 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-6xl items-center gap-4 px-5 sm:px-6">
        <a
          href="#top"
          className="rounded-lg outline-none focus-visible:ring-2 focus-visible:ring-ring/60"
          aria-label="NeoHM"
        >
          <Logo />
        </a>

        <nav
          className="ml-5 hidden items-center gap-6 md:flex"
          aria-label={t.ui.primaryNav}
        >
          {links.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className="relative text-sm text-muted-foreground transition-colors after:absolute after:-bottom-1.5 after:left-0 after:h-px after:w-0 after:bg-brand-blue after:transition-all after:duration-300 hover:text-foreground hover:after:w-full"
            >
              {link.label}
            </a>
          ))}
        </nav>

        <div className="ml-auto flex items-center gap-2">
          <LanguageSwitch />
          <a
            href={REPO_URL}
            target="_blank"
            rel="noreferrer"
            aria-label={t.nav.github}
            className={cn(
              buttonVariants({ variant: "ghost", size: "icon" }),
              "rounded-full text-muted-foreground hover:text-foreground"
            )}
          >
            <GithubIcon className="size-4" />
          </a>
          <a
            href={DOWNLOAD_URL}
            target="_blank"
            rel="noreferrer"
            className={cn(
              buttonVariants({ variant: "secondary", size: "sm" }),
              "hidden rounded-full pr-2.5 sm:inline-flex"
            )}
          >
            {t.nav.download}
            <ArrowRight data-icon="inline-end" />
          </a>
        </div>
      </div>
    </header>
  )
}
