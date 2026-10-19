import { Check } from "lucide-react"

import { MenuBar } from "@/components/site/menu-bar"
import { Reveal } from "@/components/site/reveal"
import { SectionHeading } from "@/components/site/section-heading"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useLanguage } from "@/lib/language"
import { cn } from "@/lib/utils"

export function Showcase() {
  const { t, screens } = useLanguage()
  const tabs = t.showcase.tabs

  const items = [
    {
      id: "home",
      copy: tabs.home,
      images: [{ src: screens.home, alt: t.alts.home, label: null }],
    },
    {
      id: "climate",
      copy: tabs.climate,
      images: [{ src: screens.climate, alt: t.alts.climate, label: null }],
    },
    {
      id: "devices",
      copy: tabs.devices,
      images: [
        { src: screens.light, alt: t.alts.light, label: tabs.devices.shots[0] },
        { src: screens.cover, alt: t.alts.cover, label: tabs.devices.shots[1] },
        { src: screens.media, alt: t.alts.media, label: tabs.devices.shots[2] },
      ],
    },
    {
      id: "settings",
      copy: tabs.settings,
      images: [{ src: screens.settings, alt: t.alts.settings, label: null }],
    },
  ]

  return (
    <section
      id="screens"
      className="relative scroll-mt-24 border-y border-white/[0.05] bg-white/[0.012] py-24 sm:py-28"
    >
      <div className="mx-auto max-w-6xl px-5 sm:px-6">
        <Reveal>
          <SectionHeading
            eyebrow={t.showcase.eyebrow}
            title={t.showcase.title}
            body={t.showcase.body}
          />
        </Reveal>

        <Reveal delay={80}>
          <Tabs defaultValue="home" className="mt-12 gap-8">
            <TabsList className="h-auto w-fit flex-wrap justify-start gap-1 rounded-full border border-white/[0.07] bg-white/[0.03] p-1">
              {items.map((item) => (
                <TabsTrigger
                  key={item.id}
                  value={item.id}
                  className="h-8 flex-none rounded-full px-4 text-sm text-muted-foreground data-active:bg-white/10 data-active:text-foreground dark:data-active:bg-white/10 dark:data-active:text-foreground"
                >
                  {item.copy.label}
                </TabsTrigger>
              ))}
            </TabsList>

            {items.map((item) => (
              <TabsContent
                key={item.id}
                value={item.id}
                className="outline-none data-open:animate-in data-open:duration-300 data-open:fade-in-0"
              >
                <div className="grid items-center gap-10 lg:grid-cols-[0.9fr_1.1fr] lg:gap-14">
                  <div className="flex max-w-lg flex-col gap-5">
                    <h3 className="text-2xl font-semibold tracking-[-0.02em] text-balance sm:text-3xl">
                      {item.copy.title}
                    </h3>
                    <p className="text-[15px] leading-relaxed text-muted-foreground">
                      {item.copy.body}
                    </p>
                    <ul className="mt-1 flex flex-col gap-3">
                      {item.copy.points.map((point) => (
                        <li
                          key={point}
                          className="flex items-start gap-2.5 text-sm text-foreground/85"
                        >
                          <span className="mt-0.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-brand-mint/10 text-brand-mint ring-1 ring-brand-mint/15">
                            <Check className="size-3" />
                          </span>
                          {point}
                        </li>
                      ))}
                    </ul>
                  </div>

                  <div
                    className={cn(
                      "relative flex justify-center",
                      item.images.length > 1 && "items-start"
                    )}
                  >
                    <div
                      aria-hidden="true"
                      className="pointer-events-none absolute inset-x-4 top-8 bottom-8 -z-0 rounded-[40%] bg-[radial-gradient(circle,rgba(75,126,255,0.16),transparent_68%)] blur-2xl"
                    />

                    {item.images.length > 1 ? (
                      <div className="relative z-10 grid w-full grid-cols-3 gap-2 sm:gap-3">
                        {item.images.map((image) => (
                          <figure
                            key={image.alt}
                            className="flex min-w-0 flex-col items-center gap-2"
                          >
                            <img
                              src={image.src}
                              alt={image.alt}
                              width={776}
                              height={1144}
                              loading="lazy"
                              decoding="async"
                              className="w-full drop-shadow-[0_20px_36px_rgba(1,5,14,0.55)] select-none"
                            />
                            <figcaption className="text-[11px] text-muted-foreground">
                              {image.label}
                            </figcaption>
                          </figure>
                        ))}
                      </div>
                    ) : (
                      <div className="relative z-10 w-full max-w-[388px]">
                        <MenuBar className="relative z-10" />
                        <img
                          src={item.images[0]?.src}
                          alt={item.images[0]?.alt}
                          width={776}
                          height={1144}
                          loading="lazy"
                          decoding="async"
                          className="relative z-0 block w-full drop-shadow-[0_28px_50px_rgba(1,5,14,0.55)] select-none"
                        />
                      </div>
                    )}
                  </div>
                </div>
              </TabsContent>
            ))}
          </Tabs>
        </Reveal>
      </div>
    </section>
  )
}
