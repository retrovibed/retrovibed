// developer mode is a in memory representation
// used by developers to access new/in progress/debugging
// functionality. its not meant to be used by general users.
class DeveloperMode {
  bool alpha; // alpha functionality.
  bool beta; // beta functionality.
  bool networking; // enable networking functionality ux.
  bool subscription; // force enable subscription management ux.
  bool recommendations; // toggle recommendations panel.
  bool releases; // toggle releases panel.
  bool debug; // enable debug-only ux and tuning panels.

  DeveloperMode({
    this.alpha = false,
    this.beta = false,
    this.networking = false,
    this.subscription = false,
    this.recommendations = false,
    this.releases = false,
    this.debug = false,
  });

  DeveloperMode copyWith({
    bool? alpha,
    bool? beta,
    bool? networking,
    bool? subscription,
    bool? recommendations,
    bool? releases,
    bool? debug,
  }) {
    return DeveloperMode(
      alpha: alpha ?? this.alpha,
      beta: beta ?? this.beta,
      networking: networking ?? this.networking,
      subscription: subscription ?? this.subscription,
      recommendations: recommendations ?? this.recommendations,
      releases: releases ?? this.releases,
      debug: debug ?? this.debug,
    );
  }
}
