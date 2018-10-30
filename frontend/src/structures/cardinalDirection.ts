export default function getCardinalDirection(heading: number) {
    if (heading >= 22.5 && heading < 67.5) {
      return 'northeast';
    } else if (heading >= 67.5 && heading < 112.5) {
      return 'east';
    } else if (heading >= 112.5 && heading < 157.5) {
      return 'southeast';
    } else if (heading >= 157.5 && heading < 202.5) {
      return 'south';
    } else if (heading >= 202.5 && heading < 247.5) {
      return 'southwest';
    } else if (heading >= 247.5 && heading < 292.5) {
      return 'west';
    } else if (heading >= 292.5 && heading < 337.5) {
      return 'northwest';
    }
    return 'north';
  }
