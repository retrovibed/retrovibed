import 'package:flutter/material.dart';
import 'theme.defaults.dart';

class Accordion extends StatefulWidget {
  final Widget description;
  final Widget content;
  final Widget? disabled;
  final bool expanded;

  const Accordion({
    Key? key,
    required this.description,
    required this.content,
    this.disabled,
    this.expanded = false,
  }) : super(key: key);

  @override
  State<Accordion> createState() => _AccordionState(!this.expanded);
}

class _AccordionState extends State<Accordion> {
  bool hidden;

  _AccordionState(this.hidden);

  void toggle() {
    setState(() {
      hidden = !hidden;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final borderRadius = Defaults.of(context).borderRadius;
    final isDisabled = widget.disabled != null;
    final onPressed = isDisabled ? null : toggle;
    final opacity = isDisabled ? 0.2 : 1.0;
    final icon = widget.disabled ?? Icon(hidden ? Icons.arrow_drop_up : Icons.arrow_drop_down);
    final cursor = isDisabled ? SystemMouseCursors.forbidden : SystemMouseCursors.click;

    final content = hidden
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
          child: Material(
            color: Colors.transparent,
            child: ListTile(
              shape: RoundedRectangleBorder(borderRadius: borderRadius),
              hoverColor: theme.hoverColor,
              mouseCursor: cursor,
              onTap: onPressed,
              title: widget.description,
              trailing: icon,
            ),
          ),
        ),
        Visibility(visible: !hidden, child: content),
      ],
    );
  }
}
