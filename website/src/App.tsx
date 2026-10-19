import { CallToAction } from "@/components/site/call-to-action"
import { Faq } from "@/components/site/faq"
import { Features } from "@/components/site/features"
import { Footer } from "@/components/site/footer"
import { Hero } from "@/components/site/hero"
import { Languages } from "@/components/site/languages"
import { Nav } from "@/components/site/nav"
import { Setup } from "@/components/site/setup"
import { Showcase } from "@/components/site/showcase"
import { LanguageProvider, useLanguage } from "@/lib/language"

function PageBackground() {
  return (
    <div
      aria-hidden="true"
      className="pointer-events-none absolute inset-x-0 top-0 -z-10 h-[1600px] overflow-hidden"
    >
      <div className="absolute inset-0 bg-grid-faint [mask-image:radial-gradient(ellipse_75%_55%_at_50%_8%,black,transparent)]" />
      <div className="absolute top-[-360px] left-1/2 h-[840px] w-[1120px] -translate-x-1/2 rounded-full bg-[radial-gradient(ellipse_at_center,rgba(54,101,214,0.26),transparent_64%)] blur-3xl" />
      <div className="absolute top-[560px] right-[-14%] size-[460px] rounded-full bg-[radial-gradient(circle,rgba(246,195,126,0.09),transparent_66%)] blur-3xl" />
      <div className="absolute top-[420px] left-[-12%] size-[420px] rounded-full bg-[radial-gradient(circle,rgba(183,156,255,0.09),transparent_66%)] blur-3xl" />
      <div className="absolute inset-x-0 bottom-0 h-72 bg-gradient-to-b from-transparent to-background" />
    </div>
  )
}

function Site() {
  const { t } = useLanguage()

  return (
    <div id="top" className="relative isolate min-h-svh overflow-x-clip">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:fixed focus:top-4 focus:left-4 focus:z-[100] focus:rounded-full focus:bg-white focus:px-4 focus:py-2 focus:text-sm focus:font-medium focus:text-black"
      >
        {t.ui.skip}
      </a>

      <PageBackground />
      <Nav />

      <main id="main" tabIndex={-1}>
        <Hero />
        <Features />
        <Showcase />
        <Languages />
        <Setup />
        <Faq />
        <CallToAction />
      </main>

      <Footer />
    </div>
  )
}

export function App() {
  return (
    <LanguageProvider>
      <Site />
    </LanguageProvider>
  )
}
