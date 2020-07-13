import {StoreState} from '@/StoreState';
const SunCalc = require('suncalc');

/**
 * == How dark theme works in this app ==
 *
 * The current theme mode is saved in the global app settings. `state.settings.darkThemeMode`. This mode defines how
 * the dark theme should be applied. The settings page has a toggle that turns dark theme on or off, and the user can
 * pick which specific mode they want. The settings logic is in `settings.vue` and `store.ts`.
 *
 * The actual theme is implemented via CSS variables and a separate style for dark mode. Elements have their color set
 * not via a direct color, but by referencing a CSS variable. For example, most text elements have their color as
 * `color: var(--color-fg-normal)`. In the normal theme, the fg-normal color is black. When the dark theme is enabled,
 * this variable changes to white. The CSS logic for dark mode is in `theme.ts`.
 *
 * Top-level apps have a `data-theme=` attribute on their root element. This attribute controls the current theme of the
 * app. The value of this attribute is controlled by a computed property that calls {@link DarkTheme#getCurrentCSSThemeAttribute()}.
 * When the theme is updated from the settings menu, Vue magic happens and this computed property updates to the right
 * value according to the settings, enabling or disabling the app dark theme.
 */

/** Describes a dark theme mode. */
export class DarkThemeMode {
  public static readonly OFF = new DarkThemeMode('off', 'Off', (state) => false);
  public static readonly ON = new DarkThemeMode('on', 'Always', (state) => true);
  public static readonly AT_NIGHT = new DarkThemeMode('timeauto', 'Auto',
    (state) => {
      const times = SunCalc.getTimes(state.now, 42.7299, -73.6766, 72 /* lat, long, and elevation(m) of RPI Union. */);
      return state.now < times.sunrise || state.now > times.sunset;
    },
  );

  /**
   * Every mode defined in this class. Computed dynamically using reflection.
   */
  public static allModes(): DarkThemeMode[] {
    return Object.values(DarkThemeMode)
      .filter((o) => o instanceof DarkThemeMode);
  }

  /** Converts a string ID to a {@link DarkThemeMode} object. */
  public static fromString(darkThemeId: string): DarkThemeMode {
    // Reflection is used to find every defined DarkThemeMode.
    const mode = DarkThemeMode.allModes().find((m) => m.id === darkThemeId);
    if (mode === undefined) {
      throw new TypeError('Invalid ID :(');
    } else {
      return mode;
    }
  }

  /** The internal ID of this mode. This is how the mode is saved in storage. Should not be edited on existing modes. */
  public readonly id: string;
  /** The name of this mode when seen by the user. */
  public description: string;
  /** Asks if dark theme should be enabled at this moment. Some modes have dynamic behavior. */
  public isDarkThemeVisible: (state: StoreState) => boolean;

  constructor(id: string, description: string, visibleFn: (state: StoreState) => boolean) {
    this.id = id;
    this.description = description;
    this.isDarkThemeVisible = visibleFn;
  }
}

/** Holds all valid instances of {@link DarkThemeMode}. */
export class DarkTheme {
  /**
   * Determines if the dark theme should be displayed to the user at this time. Anywhere in the site that needs to check
   * the status of dark theme should use this function.
   *
   * This is DIFFERENT from checking the mode in {@link StoreState} because some modes are dynamic.
   */
  public static isDarkThemeVisible(state: StoreState): boolean {
    return DarkThemeMode.fromString(state.settings.darkThemeMode).isDarkThemeVisible(state);
  }

  /**
   * Computes the appropriate value of the `data-theme` attribute for the currently active theme. Every root
   * page on the site that is to have a dark theme needs to bind the result of this function to its <body> tag.
   *
   * Example:
   *   computed: {
   *     ...
   *     currentCSSTheme(): string {
   *       return DarkTheme.getCurrentCSSThemeAttribute(this.$store.state);
   *     },
   *   },
   */
  public static getCurrentCSSThemeAttribute(state: StoreState): string {
    return DarkTheme.isDarkThemeVisible(state) ? 'dark' : '';
  }

  private constructor() {}
}
