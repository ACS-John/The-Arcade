Aphelion picture-field art

Aphelion.panda (in the folder above this one) shows a picture in a 16-row by
40-column PICTURE field near the top-left of its screen, referenced by path
from Aphelion\art\. The field resizes whatever you drop in here with
:ISOTROPIC (aspect ratio kept, no stretching), so a roughly 5:2 landscape
image (e.g. 800x320, or any similar ratio) will fit best. Supported formats
include JPG, PNG, BMP, GIF, and a few others -- see
context/br_tree/20-io-screen/controls/Picture.md.

If a file below doesn't exist, the game just skips the picture for that
scene (checked with EXISTS() before drawing) -- nothing breaks either way.

Expected filenames, one per scene group:
  title.jpg        - title screen
  crashsite.jpg     - waking up / examining the suit / scanning / calling out
  wakehub.jpg       - ECHO comes online, choosing tower/forest/wait
  tower.jpg         - the signal tower branch
  forest.jpg        - the glowing forest branch
  ending-dusk.jpg   - ending: Swallowed by Dusk
  ending-loop.jpg   - ending: The Loop
  ending-signal.jpg - ending: The Signal Answered
  ending-roots.jpg  - ending: Roots
  ending-ship.jpg   - ending: The Ship Remembers

The complete illustrated set is installed at 1983x793 pixels per image.
Open gallery.html to review all scenes (includes ending spoilers).
See ART-DIRECTION.md for continuity notes and prompts.json for generation prompts.
