import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import './screens.dart' as screens;

NodeState? of(BuildContext context) {
  return Node.of(context);
}

Future<T> asyncfn<T>(
  BuildContext context,
  Widget Function(Completer<T> completion) builder,
) {
  final completion = Completer<T>();
  final node = of(context);
  node?.push(
    _DismissGuard(
      onDismiss: () {
        if (completion.isCompleted) return;
        // Content was removed (e.g. tapping outside the modal) without an
        // explicit completion value; only succeeds for nullable/void T.
        try {
          completion.complete();
        } catch (e) {
          print("failed to dismiss async modal ${e}");
        }
      },
      child: builder(completion),
    ),
  );
  return completion.future.whenComplete(() => node?.push(null));
}

class _DismissGuard extends StatefulWidget {
  final Widget child;
  final VoidCallback onDismiss;

  const _DismissGuard({required this.child, required this.onDismiss});

  @override
  State<_DismissGuard> createState() => _DismissGuardState();
}

class _DismissGuardState extends State<_DismissGuard> {
  @override
  Widget build(BuildContext context) => widget.child;

  @override
  void dispose() {
    widget.onDismiss();
    super.dispose();
  }
}

class Node extends StatefulWidget {
  final Widget child;
  final AlignmentGeometry alignment;

  const Node(this.child, {super.key, this.alignment = Alignment.center});

  static NodeState? of(BuildContext context) {
    return context.findAncestorStateOfType<NodeState>();
  }

  @override
  State<StatefulWidget> createState() => NodeState();
}

class NodeState extends State<Node> {
  static const zeromodal = const SizedBox();
  final FocusScopeNode _selffocus = FocusScopeNode(debugLabel: "modal.node");
  Widget current = zeromodal;
  List<Widget> stack = [];

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void push(Widget? m) {
    setState(() {
      _selffocus.requestFocus();
      if (m == null) {
        setState(() {
          current = stack.isEmpty ? zeromodal : stack.last;
          stack.remove(current);
        });
        return;
      }

      setState(() {
        if (current == NodeState.zeromodal) {
          current = m;
          return;
        }

        stack = stack + [current];
        current = m;
      });
    });
  }

  void reset() {
    setState(() {
      current = zeromodal;
    });
  }

  @override
  Widget build(BuildContext context) {
    return screens.Overlay.tappable(
      widget.child,
      overlay: LayoutBuilder(
        builder: (context, constraints) {
          return FocusScope(
            canRequestFocus: true,
            autofocus: true,
            node: _selffocus,
            onKeyEvent: (node, event) {
              if (event is KeyDownEvent) {
                return KeyEventResult.ignored;
              }
              if (event.logicalKey != LogicalKeyboardKey.escape || (stack.isEmpty && current == NodeState.zeromodal)) {
                return KeyEventResult.ignored;
              }

              push(null);
              return KeyEventResult.handled;
            },
            child: Visibility(
              visible: current != zeromodal,
              child: screens.Masked(
                Center(child: SingleChildScrollView(child: current)),
                reset: reset,
              ),
            ),
          );
        },
      ),
      alignment: widget.alignment,
      onTap: current != zeromodal ? reset : null,
    );
  }
}
