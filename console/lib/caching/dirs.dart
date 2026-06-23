class DirsWellKnown {
  String cache = '';
  DirsWellKnown({
    this.cache = '',
  });

  factory DirsWellKnown.xdg(({String cacheDir, String configDir, String dataDir, String downloadDir}) xdg) {
    return DirsWellKnown(
      cache: xdg.cacheDir,
    );
  }
}

var _global = DirsWellKnown();

void setglobal(DirsWellKnown d) {
  _global = d;
}

DirsWellKnown global() {
  return _global;
}
