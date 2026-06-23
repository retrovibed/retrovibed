import 'dart:io';
import 'package:path/path.dart' as p;
import 'package:flutter/material.dart';
import 'package:retrovibed/caching.dart' as fscache;
import 'guarded.dart';

// Disclaimer is for showing a notification a single time.
class Disclaimer extends StatefulWidget {
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

  static bool disclaimerpath(String id) => defaultcache('disclaimer', id);
  static bool defaultcache(String prefix, String id) =>
      fscache.Dir(Directory(p.join(fscache.global().cache, prefix))).maybe(id, () => false);

  static void acknowledge(String id) => defaultwrite('disclaimer', id);
  static void defaultwrite(String prefix, String id) =>
      fscache.Dir(Directory(p.join(fscache.global().cache, prefix))).write(id, true);

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
