# Image credits

Both source images are **public domain** by age, and both scans are released
without restrictions by the institutions that hold them.

Attribution isn't legally required for public-domain works. It's here anyway,
because being able to say exactly where an asset came from matters.

---

## `hestia-hero.png`

**Standbeeld van Hestia** (the *Hestia Giustiniani*)
Rijksmuseum, Amsterdam, object `RP-P-2015-3-16`. Plate from the
*Galleria Giustiniana del Marchese Vincenzo Giustiniani*, Rome, 1636–1647.
Signed on the base *C. Mellan Gall. del. et f.* (Claude Mellan). Engraving on
paper. Source scan 4308 × 5592 px, **CC0**, via
[Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Standbeeld_van_Hestia_Galleria_Giustiniana_del_Marchese_Vincenzo_Giustiniani_I_(serietitel),_RP-P-2015-3-16.jpg).

The statue is a Roman copy of a Greek original of about 470 BC, now in the
Museo Torlonia, Rome. It's one of the very few surviving ancient statues of
Hestia: veiled, plainly dressed and standing still, as the goddess who never
leaves the hearth.

## `hestia-hearth.jpg`

**Vesta**, from Hollar's series of the Greek gods.
Wenceslaus Hollar (1607–1677). Etching. Scan from the
University of Toronto Wenceslaus Hollar Digital Collection, **public domain**,
via
[Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Wenceslas_Hollar_-_The_Greek_gods._Vesta.jpg).

Vesta is Hestia's Roman counterpart. She sits on the left; the fire on the
altar is the hearth itself.

---

## Processing

Both files are derivatives, not the raw scans.

**Masthead figure:**

1. Cropped to the figure and its plinth.
2. Figure isolated from the background with a segmentation model (`rembg`,
   `birefnet-general`; `isnet-general-use` and `u2net` were tried and left
   holes in the hatching or soft edges).
3. Mask applied at **full scan resolution**, so the burin hatching survives.
4. Levels normalised from percentiles of the masked region, then floated on
   pure black.

**Hearth plate:** cropped to the printed border by detecting the inked
bounding box, converted to greyscale and normalised. The figure isn't cut out
because the altar flame overlaps the background hatching too closely for a
clean mask.

**Masthead type and colour:** Roboto Mono throughout, on `#000000`, with white
line art and type. `#F97316` (ember orange) is used exactly once, for the short
rule: the hearth fire.
