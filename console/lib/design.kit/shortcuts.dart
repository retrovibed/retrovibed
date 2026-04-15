import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/design.kit/help.dart';

typedef ShortcutBinding = (Widget help, KeyEventResult Function() handler);

// Like Shortcuts but callbacks return bool to indicate whether the
// event was handled. Returning false allows the event to continue propagating.
// Bindings may include a help widget that is registered with the nearest
// HelpScope ancestor for display in help overlays.
class Shortcuts extends StatefulWidget {
  final Map<ShortcutActivator, ShortcutBinding> bindings;
  final Widget child;

  const Shortcuts(this.child, {super.key, required this.bindings});

  @override
  State<Shortcuts> createState() => _ShortcutsState();
}

class _ShortcutsState extends State<Shortcuts> {
  HelpScopeState? _helpScope;
  List<Widget> _registered = [];

  void _syncHelp() {
    final scope = HelpScope.of(context);
    for (final w in _registered) {
      _helpScope?.unregisterGlobal(w);
    }
    _helpScope = scope;
    _registered =
        widget.bindings.entries.map((e) {
          final activator = e.key;
          String labelText = '';

          if (activator is SingleActivator) {
            // Formats as "Ctrl+Shift+S" etc.
            final modifiers = [
              if (activator.control) 'Ctrl',
              if (activator.alt) 'Alt',
              if (activator.shift) 'Shift',
              if (activator.meta) 'Meta',
              activator.trigger.keyLabel,
            ];
            labelText = modifiers.join('+');
          }

          return Hint(label: Text(labelText), description: e.value.$1);
        }).toList();
    for (final w in _registered) {
      scope?.registerGlobal(w);
    }
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    _syncHelp();
  }

  @override
  void didUpdateWidget(Shortcuts oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.bindings != oldWidget.bindings) {
      _syncHelp();
    }
  }

  @override
  void dispose() {
    for (final w in _registered) {
      _helpScope?.unregisterGlobal(w);
    }
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Focus(
      canRequestFocus: false,
      onKeyEvent: (node, event) {
        for (final entry in widget.bindings.entries) {
          if (entry.key.accepts(event, HardwareKeyboard.instance)) {
            return entry.value.$2();
          }
        }
        return KeyEventResult.ignored;
      },
      child: widget.child,
    );
  }
}
