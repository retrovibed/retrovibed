import 'dart:io';
import 'package:flutter/material.dart';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';
import 'package:retrovibed/caching/fscache.dart' as fscache;
import 'flutterx.dart';
import 'screens.dart' as screens;
import 'buttons.dart';
import 'help.dart' show HelpScope;
import 'container.dart' as _c;
import 'theme.defaults.dart';

/// Wraps [child] and displays [content] as a one-time masked overlay the
/// first time this widget is mounted. The [cacheid] key is persisted to disk
/// so the overlay is never shown again after the first dismissal.
class HelpAuto extends StatefulWidget {
  final Widget child;
  final Widget title;
  final Widget content;
  final String cacheid;
  const HelpAuto(this.child, {required this.title, required this.content, required this.cacheid, super.key});

  @override
  State<HelpAuto> createState() => _HelpAutoState();
}

class _HelpAutoState extends State<HelpAuto> {
  bool _visible = false;

  @override
  void initState() {
    super.initState();
    _activate();
  }

  Future<void> _activate() async {
    final cacheDir = await getApplicationCacheDirectory();
    final cache = fscache.Dir(Directory(p.join(cacheDir.path, 'help')));
    final alreadyActivated = cache.maybe<bool>(widget.cacheid, () => false);
    if (!alreadyActivated) {
      cache.write<bool>(widget.cacheid, true);
      postframe(() {
        if (mounted) setState(() => _visible = true);
      });
    }
  }

  void _close() => setState(() => _visible = false);

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    final activate = defaults.mobile
        ? "Shake the phone to activate/deactivate help"
        : "Press Alt+? at any time to activate/deactivate help overlay";
    return screens.Overlay(
      widget.child,
      overlay: Visibility(
        visible: _visible,
        child: screens.Masked(
          Center(
            child: SingleChildScrollView(
              child: _c.Container(
                padding: defaults.padding,
                margin: defaults.margin,
                constraints: const BoxConstraints(maxWidth: 512),
                Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  spacing: defaults.spacing,
                  children: [
                    Row(
                      children: [
                        widget.title,
                        Spacer(),
                        buttons.remove(onPressed: _close),
                      ],
                    ),
                    const Divider(),
                    widget.content,
                    Text(activate),
                    ...HelpScope.of(context)?.globals ?? [],
                  ],
                ),
              ),
            ),
          ),
          reset: _close,
        ),
      ),
    );
  }
}
