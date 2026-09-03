import 'dart:io';
import 'package:path/path.dart' as p;
import 'package:flutter/material.dart';
import 'package:retrovibed/caching.dart' as fscache;
import 'guarded.dart';

// Disclaimer is for showing a notification a single time.
class Disclaimer extends StatefulWidget {
  static final _prefix = 'disclaimer';
  final Widget child;
  final Widget overlay;
  final String cacheid;
  final bool Function(String id) cached;
  const Disclaimer(
    this.child, {
    super.key,
    required this.cacheid,
    required this.overlay,
    this.cached = Disclaimer.disclaimerpath,
  });

  static fscache.Dir defaultcachedir() {
    return fscache.Dir(Directory(p.join(fscache.global().cache, _prefix)));
  }

  static bool disclaimerpath(String id) => defaultcache(_prefix, id);
  static bool defaultcache(String prefix, String id) =>
      fscache.Dir(Directory(p.join(fscache.global().cache, prefix))).maybe(id, () => false);

  static void acknowledge(String id) => defaultwrite(_prefix, id);
  static void defaultwrite(String prefix, String id) {
    fscache.Dir(Directory(p.join(fscache.global().cache, prefix))).write(id, true);
  }

  @override
  State<Disclaimer> createState() => _DisclaimerState();
}

class _DisclaimerState extends State<Disclaimer> {
  @override
  Widget build(BuildContext context) {
    return Guarded(
      child: widget.child,
      overlay: widget.overlay,
      enabled: !widget.cached(widget.cacheid),
    );
  }
}
