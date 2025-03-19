import 'package:flutter/material.dart';

class Accordion extends StatefulWidget {
  final Widget description;
  final Widget content;
  final Widget? disabled;

  const Accordion({
    Key? key,
    required this.description,
    required this.content,
    this.disabled,
  }) : super(key: key);

  @override
  State<Accordion> createState() => _AccordionState();
}

class _AccordionState extends State<Accordion> {
  bool hidden = true;

  void toggle() {
    setState(() {
      hidden = !hidden;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDisabled = widget.disabled != null;
    final onPressed = isDisabled ? null : toggle;
    final opacity = isDisabled ? 0.2 : 1.0;
    final icon =
        widget.disabled ??
        Icon(hidden ? Icons.arrow_drop_up : Icons.arrow_drop_down);
    final cursor =
        isDisabled ? SystemMouseCursors.forbidden : SystemMouseCursors.click;

    final content =
        hidden
            ? Container()
            : Container(
              padding: theme.buttonTheme.padding,
              child: widget.content,
            );

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Opacity(
          opacity: opacity,
          child: ListTile(
            hoverColor: theme.hoverColor,
            mouseCursor: cursor,
            onTap: onPressed,
            title: widget.description,
            trailing: icon,
          ),
        ),
        content,
      ],
    );
  }
}
