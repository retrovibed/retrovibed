// developer mode is a in memory representation
// used by developers to access new/in progress/debugging
// functionality. its not meant to be used by general users.
class DeveloperMode {
  bool alpha = false; // alpha functionality.
  bool beta = false; // beta functionality.
  bool networking = false; // enable networking functionality ux.
  bool subscription = false; // force enable subscription management ux.

  DeveloperMode();
}
