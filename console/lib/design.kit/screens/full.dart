import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:window_manager/window_manager.dart';
import 'package:retrovibed/env.dart' as env;
import 'package:retrovibed/windowx.dart';

class Full extends StatefulWidget {
  final Widget? child;
  final WindowManagerX windowManager;
  Full(this.child, {super.key, WindowManagerX? windowManager}) : windowManager = windowManager ?? windowx;

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

class _FullState extends State<Full> with WindowListener {
  bool chromeless = false;

  void toggle() {
    final next = !chromeless;
    setState(() => chromeless = next);

    if (env.boolean(
      env.vars.WindowManagerNativeFullScreen,
      fallback: Platform.isLinux || Platform.isWindows || Platform.isMacOS,
    )) {
      widget.windowManager.setFullScreen(next).catchError((cause) {
        print("failed to toggle fullscreen mode ${next} ${cause}");
      });
    } else {
      SystemChrome.setEnabledSystemUIMode(
        next ? SystemUiMode.immersiveSticky : SystemUiMode.edgeToEdge,
      );
    }
  }

  @override
  void initState() {
    super.initState();
    widget.windowManager.addListener(this);
    widget.windowManager.isFullScreen().then((kFullscreen) {
      setState(() {
        chromeless = kFullscreen;
      });
    }).catchError((cause) {
      print("failed to read initial fullscreen state ${cause}");
    });
  }

  @override
  void dispose() {
    widget.windowManager.removeListener(this);
    super.dispose();
  }

  // Keeps chromeless in sync with fullscreen changes that don't go through
  // toggle() - an OS-level shortcut, another entry point, or a setFullScreen
  // call above that failed silently (its catchError only logs).
  @override
  void onWindowEnterFullScreen() => setState(() => chromeless = true);

  @override
  void onWindowLeaveFullScreen() => setState(() => chromeless = false);

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
