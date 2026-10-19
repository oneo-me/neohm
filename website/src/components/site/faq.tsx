import { ArrowUpRight } from "lucide-react"

import { Reveal } from "@/components/site/reveal"
import { SectionHeading } from "@/components/site/section-heading"
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"
import { useLanguage } from "@/lib/language"
import { ISSUES_URL } from "@/lib/links"

export function Faq() {
  const { t } = useLanguage()

  return (
    <section id="faq" className="scroll-mt-24 py-24 sm:py-28">
      <div className="mx-auto grid max-w-6xl gap-12 px-5 sm:px-6 lg:grid-cols-[0.82fr_1.18fr] lg:gap-16">
        <Reveal>
          <div className="lg:sticky lg:top-28">
            <SectionHeading
              eyebrow={t.faq.eyebrow}
              title={t.faq.title}
              body={t.faq.body}
            />
            <a
              href={ISSUES_URL}
              target="_blank"
              rel="noreferrer"
              className="mt-6 inline-flex items-center gap-1.5 text-sm text-brand-blue transition-colors hover:text-foreground"
            >
              {t.faq.more} {t.faq.link}
              <ArrowUpRight className="size-3.5" />
            </a>
          </div>
        </Reveal>

        <Reveal delay={80}>
          <Accordion
            defaultValue={["faq-0"]}
            className="rounded-2xl border border-white/[0.07] bg-card/40 px-5 backdrop-blur-sm"
          >
            {t.faq.items.map((item, index) => (
              <AccordionItem
                key={item.question}
                value={`faq-${index}`}
                className="not-last:border-b-white/[0.06]"
              >
                <AccordionTrigger className="py-4 text-[15px] font-medium hover:no-underline">
                  {item.question}
                </AccordionTrigger>
                <AccordionContent className="pb-4 text-[13.5px] leading-relaxed text-muted-foreground">
                  {item.answer}
                </AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </Reveal>
      </div>
    </section>
  )
}
