import { ArrowUpRight } from "lucide-react"

import { Logo } from "@/components/site/logo"
import { Separator } from "@/components/ui/separator"
import { useLanguage } from "@/lib/language"
import { REPO_URL } from "@/lib/links"

const links = [
  { key: "github" as const, href: REPO_URL },
  { key: "readme" as const, href: `${REPO_URL}#readme` },
  {
    key: "architecture" as const,
    href: `${REPO_URL}/blob/main/ARCHITECTURE.md`,
  },
  {
    key: "changelog" as const,
    href: `${REPO_URL}/blob/main/CHANGELOG.md`,
  },
]

export function Footer() {
  const { t } = useLanguage()
  const year = new Date().getFullYear()

  return (
    <footer className="border-t border-white/[0.06] bg-black/10">
      <div className="mx-auto max-w-6xl px-5 py-14 sm:px-6">
        <div className="flex flex-col gap-10 md:flex-row md:items-start md:justify-between">
          <div className="max-w-sm">
            <Logo />
            <p className="mt-4 text-sm text-foreground/80">
              {t.footer.tagline}
            </p>
            <p className="mt-2 text-xs leading-relaxed text-muted-foreground/85">
              {t.footer.note}
            </p>
          </div>

          <div>
            <h3 className="text-[10px] font-medium tracking-[0.16em] text-muted-foreground uppercase">
              {t.footer.project}
            </h3>
            <ul className="mt-4 grid grid-cols-2 gap-x-10 gap-y-2.5">
              {links.map((link) => (
                <li key={link.key}>
                  <a
                    href={link.href}
                    target="_blank"
                    rel="noreferrer"
                    className="group inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {t.footer.links[link.key]}
                    <ArrowUpRight className="size-3 opacity-0 transition-opacity group-hover:opacity-60" />
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <Separator className="my-9 bg-white/[0.06]" />

        <div className="flex flex-col gap-2 text-xs text-muted-foreground/80 sm:flex-row sm:items-center sm:justify-between">
          <span>
            © {year} NeoHM · {t.footer.builtWith}
          </span>
          <span>{t.footer.madeBy}</span>
        </div>
      </div>
    </footer>
  )
}
