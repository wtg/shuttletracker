import {StoreState} from '@/StoreState';

/** Describes a dark theme mode. */
export class DarkThemeMode {
  /** The internal ID of this mode. This is how the mode is saved in storage. Should not be edited on existing modes. */
  public readonly id: string;
  /** The name of this mode when seen by the user. */
  public description: string;
  /** Asks if dark theme should be enabled at this moment. Some modes have dynamic behavior. */
  public isDarkThemeVisible: () => boolean;

  constructor(id: string, description: string, visibleFn: () => boolean) {
    this.id = id;
    this.description = description;
    this.isDarkThemeVisible = visibleFn;
  }
}

/** Holds all valid instances of {@link DarkThemeMode}. */
export class DarkTheme {
  public static readonly OFF = new DarkThemeMode('off', 'Off', () => false);
  public static readonly ON = new DarkThemeMode('on', 'Always', () => true);
  public static readonly AT_NIGHT = new DarkThemeMode('timeauto', 'At night 🌙',
    () => {
      // Simple hour check for nighttime. Future improvement is to use a time library to get the actual sunset and
      // sunrise times for the local user.
      const hr = new Date().getHours();
      return 17 <= hr || hr <= 7;
    },
  );
  /**
   * Every mode defined in this class. Computed dynamically using reflection.
   * This MUST come after all instances of {@link DarkThemeMode}.
   */
  public static readonly allModes = Object.values(DarkTheme)
    .filter((o) => o instanceof DarkThemeMode);

  /**
   * Determines if the dark theme should be displayed to the user at this time. Anywhere in the site that needs to check
   * the status of dark theme should use this function.
   *
   * This is DIFFERENT from checking the mode in {@link StoreState} because some modes are dynamic.
   */
  public static isDarkThemeVisible(state: StoreState): boolean {
    return this.darkThemeObjFromString(state.settings.darkThemeMode).isDarkThemeVisible();
  }

  /** Converts a string ID to a {@link DarkThemeMode} object. */
  public static darkThemeObjFromString(darkThemeId: string): DarkThemeMode {
    const mode = DarkTheme.allModes.find((m) => m.id === darkThemeId);
    if (mode === undefined) {
      throw new TypeError('Invalid ID :(');
    } else {
      return mode;
    }
  }

  private constructor() {}
}
