import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:window_manager/window_manager.dart';
import 'package:retrovibed/env.dart' as env;

class Full extends StatefulWidget {
  final Widget? child;
  const Full(this.child, {super.key});

  static _FullState? of(BuildContext context) {
    return context.findAncestorStateOfType<_FullState>();
  }

  static bool nochrome(BuildContext context) {
    return of(context)?.chromeless ?? false;
  }

  static void toggle(BuildContext context) {
    of(context)?.toggle();
  }

  @override
  State<Full> createState() => _FullState();
}

class _FullState extends State<Full> {
  bool chromeless = false;

  void toggle() {
    final next = !chromeless;
    setState(() => chromeless = next);

    if (env.boolean(
      env.vars.WindowManagerNativeFullScreen,
      fallback: Platform.isLinux || Platform.isWindows || Platform.isMacOS,
    )) {
      windowManager.setFullScreen(next).catchError((cause) {
        print("failed to toggle fullscreen mode ${next} ${cause}");
      });
    } else {
      SystemChrome.setEnabledSystemUIMode(
        next ? SystemUiMode.immersiveSticky : SystemUiMode.edgeToEdge,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final media = MediaQuery.of(context);
    return ConstrainedBox(
      constraints: BoxConstraints(
        minWidth: media.size.width,
        minHeight: media.size.height,
      ),
      child: widget.child,
    );
  }
}
