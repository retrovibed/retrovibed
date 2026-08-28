import 'package:flutter/material.dart';
import 'package:window_manager/window_manager.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn/login.dart' as login;
import 'package:retrovibed/windowx.dart';

class Hamburger extends StatefulWidget {
  final WindowManagerX windowManager;
  Hamburger({super.key, WindowManagerX? windowManager}) : windowManager = windowManager ?? windowx;

  @override
  State<Hamburger> createState() => _HamburgerState();
}

class _HamburgerState extends State<Hamburger> with WindowListener {
  bool _maximized = false;

  @override
  void initState() {
    super.initState();
    widget.windowManager.addListener(this);
    widget.windowManager.isMaximized().then((v) => setState(() => _maximized = v)).catchError((cause) {
      print("failed to read initial maximized state ${cause}");
    });
  }

  @override
  void dispose() {
    widget.windowManager.removeListener(this);
    super.dispose();
  }

  @override
  void onWindowMaximize() => setState(() => _maximized = true);

  @override
  void onWindowUnmaximize() => setState(() => _maximized = false);

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<String>(
      tooltip: "menu",
      position: PopupMenuPosition.under,
      color: Theme.of(context).colorScheme.surface,
      surfaceTintColor: Theme.of(context).colorScheme.surface,
      icon: Icon(Icons.menu),
      itemBuilder:
          (context) => [
            PopupMenuItem(
              value: "minmax",
              child: Row(
                spacing: 8,
                children: [
                  Icon(_maximized ? Icons.fullscreen_exit : Icons.fullscreen),
                  Text(_maximized ? "Minimize" : "Maximize"),
                ],
              ),
              onTap: () => _maximized ? widget.windowManager.unmaximize() : widget.windowManager.maximize(),
            ),
            PopupMenuItem(
              value: "help",
              child: Tooltip(
                message: "shortcut (alt+?)",
                child: Row(
                  spacing: 8,
                  children: [
                    Icon(Icons.help),
                    Text("help"),
                  ],
                ),
              ),
              onTap: () {
                ds.HelpScope.of(context)?.toggle();
              },
            ),
            PopupMenuItem(
              value: "logout",
              child: Row(
                spacing: 8,
                children: [
                  Icon(Icons.logout),
                  Text("Logout"),
                ],
              ),
              onTap: () {
                login.Login.logout(context);
              },
            ),
            PopupMenuItem(
              value: "exit",
              child: Row(
                spacing: 8,
                children: [
                  Icon(Icons.close),
                  Text("Exit"),
                ],
              ),
              onTap: () => widget.windowManager.close(),
            ),
          ],
    );
  }
}
