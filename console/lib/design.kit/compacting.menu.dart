import 'package:collection/collection.dart';
import 'package:flutter/material.dart';
import 'flutterx.dart';
import 'theme.defaults.dart';

class CompactingMenu extends StatefulWidget {
  final List<Widget> children;
  final MainAxisAlignment mainAxisAlignment;
  final MainAxisSize mainAxisSize;
  final CrossAxisAlignment crossAxisAlignment;

  final Widget icon;

  const CompactingMenu(
    this.children, {
    super.key,
    this.mainAxisAlignment = MainAxisAlignment.center,
    this.mainAxisSize = MainAxisSize.min,
    this.crossAxisAlignment = CrossAxisAlignment.center,
    this.icon = const Icon(Icons.more_vert),
  });

  /// Wraps [child] so it stays pinned in the visible row
  /// rather than being moved into the overflow menu when compact.
  static Widget pinned(Widget child, {Key? key}) => _Pinned(key: key, child: child);

  /// Wraps [child] so it stays pinned in the visible row and expands
  /// to fill the remaining space, rather than being moved into the
  /// overflow menu when compact.
  static Widget expanded(Widget child, {Key? key}) => _Expanded(key: key, child: child);

  @override
  State<CompactingMenu> createState() => _CompactingMenuState();
}

class _CompactingMenuState extends State<CompactingMenu> {
  bool _open = false;

  List<Widget> _renderRow(List<Widget> items) => [
    for (final w in items) w is _Expanded ? Expanded(child: w.child) : w,
  ];

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);

    return Builder(
      builder: (context) {
        final compact = defaults.isCompact;
        final spacing = 1.0;

        if (!compact) {
          return Column(
            mainAxisAlignment: widget.mainAxisAlignment,
            mainAxisSize: widget.mainAxisSize,
            crossAxisAlignment: widget.crossAxisAlignment,
            children: [
              Row(spacing: spacing, children: _renderRow(widget.children)),
            ],
          );
        }

        final groups = groupBy(widget.children, (w) => w is _Pinned || w is _Expanded);
        final visible = groups[true] ?? [];
        final menuItems = groups[false] ?? [];

        if (menuItems.isEmpty) {
          return Column(
            mainAxisAlignment: widget.mainAxisAlignment,
            mainAxisSize: widget.mainAxisSize,
            crossAxisAlignment: widget.crossAxisAlignment,
            children: [
              Row(mainAxisAlignment: MainAxisAlignment.center, spacing: spacing, children: _renderRow(widget.children)),
            ],
          );
        }

        return Column(
          mainAxisAlignment: widget.mainAxisAlignment,
          mainAxisSize: widget.mainAxisSize,
          crossAxisAlignment: widget.crossAxisAlignment,
          children: [
            Row(
              mainAxisSize: MainAxisSize.max,
              mainAxisAlignment: MainAxisAlignment.center,
              spacing: spacing,
              children: [
                ..._renderRow(visible),
                if (menuItems.isNotEmpty)
                  IconButton(
                    onPressed: () {
                      setState(() => _open = !_open);
                      postframe(() => Scrollable.ensureVisible(context));
                    },
                    icon: widget.icon,
                  ),
              ],
            ),
            if (_open && menuItems.isNotEmpty)
              Padding(
                padding: EdgeInsetsGeometry.symmetric(vertical: defaults.padding.vertical / 4),
                child: Row(
                  mainAxisSize: MainAxisSize.max,
                  mainAxisAlignment: MainAxisAlignment.center,
                  spacing: spacing,
                  children: menuItems,
                ),
              ),
          ],
        );
      },
    );
  }
}

/// Wraps a widget to keep it pinned in the visible row of a [CompactingMenu]
/// widget rather than being moved into the overflow menu when compact.
class _Pinned extends StatelessWidget {
  final Widget child;

  const _Pinned({super.key, required this.child});

  @override
  Widget build(BuildContext context) => child;
}

/// Wraps a widget to keep it pinned in the visible row of a [CompactingMenu]
/// and expand to fill the remaining space, rather than being moved into
/// the overflow menu when compact.
class _Expanded extends StatelessWidget {
  final Widget child;

  const _Expanded({super.key, required this.child});

  @override
  Widget build(BuildContext context) => child;
}
