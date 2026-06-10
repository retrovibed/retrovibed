import 'package:collection/collection.dart';
import 'package:flutter/material.dart';
import 'flutterx.dart';
import 'theme.defaults.dart';

class CompactingMenu extends StatefulWidget {
  final Widget child;
  final List<Widget> leading;
  final List<Widget> trailing;
  final MainAxisAlignment mainAxisAlignment;
  final MainAxisSize mainAxisSize;
  final CrossAxisAlignment crossAxisAlignment;

  final Widget icon;

  const CompactingMenu(
    this.child, {
    super.key,
    this.leading = const [],
    this.trailing = const [],
    this.mainAxisAlignment = MainAxisAlignment.start,
    this.mainAxisSize = MainAxisSize.min,
    this.crossAxisAlignment = CrossAxisAlignment.center,
    this.icon = const Icon(Icons.more_vert),
  });

  /// Wraps [child] so it stays pinned in the leading/trailing row
  /// rather than being moved into the overflow menu when compact.
  static Widget pinned(Widget child, {Key? key}) => _Pinned(key: key, child: child);

  @override
  State<CompactingMenu> createState() => _CompactingMenuState();
}

class _CompactingMenuState extends State<CompactingMenu> {
  bool _open = false;

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);

    return Builder(
      builder: (context) {
        final compact = defaults.isCompact;
        final spacing = 1.0;

        if (!compact) {
          return Row(
            spacing: spacing,
            children: [
              ...widget.leading,
              Expanded(child: widget.child),
              ...widget.trailing,
            ],
          );
        }

        final leadingGroups = groupBy(widget.leading, (w) => w is _Pinned);
        final trailingGroups = groupBy(widget.trailing, (w) => w is _Pinned);
        final pinnedLeading = leadingGroups[true] ?? [];
        final pinnedTrailing = trailingGroups[true] ?? [];
        final menuItems = [...?leadingGroups[false], ...?trailingGroups[false]];

        if (menuItems.isEmpty) {
          return Row(
            spacing: spacing,
            children: [
              ...widget.leading,
              Expanded(child: widget.child),
              ...widget.trailing,
            ],
          );
        }

        return Column(
          mainAxisAlignment: widget.mainAxisAlignment,
          mainAxisSize: widget.mainAxisSize,
          crossAxisAlignment: widget.crossAxisAlignment,
          children: [
            Row(
              spacing: spacing,
              children: [
                ...pinnedLeading,
                Expanded(child: widget.child),
                ...pinnedTrailing,
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
              Row(
                mainAxisSize: MainAxisSize.max,
                spacing: spacing,
                children: [
                  for (final w in menuItems) Expanded(child: w),
                ],
              ),
          ],
        );
      },
    );
  }
}

/// Wraps a widget to keep it pinned in the leading/trailing row of a [CompactingMenu]
/// widget rather than being moved into the overflow menu when compact.
class _Pinned extends StatelessWidget {
  final Widget child;

  const _Pinned({super.key, required this.child});

  @override
  Widget build(BuildContext context) => child;
}
