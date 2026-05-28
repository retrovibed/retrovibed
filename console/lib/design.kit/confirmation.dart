import 'dart:async';
import 'package:flutter/material.dart';
import './theme.defaults.dart';
import 'container.dart' as _container;
import 'empty.dart';

class Confirmation extends StatelessWidget {
  final Widget content;
  final Widget confirmation;
  final Widget cancellation;
  final void Function(BuildContext)? onConfirm;
  final void Function(BuildContext)? onCancel;

  const Confirmation({
    super.key,
    required this.content,
    required this.confirmation,
    required this.cancellation,
    this.onConfirm,
    this.onCancel,
  });

  factory Confirmation.ok({
    Key? key,
    required Widget content,
    void Function(BuildContext)? onConfirm,
    void Function(BuildContext)? onCancel,
  }) {
    return Confirmation(
      key: key,
      content: content,
      confirmation: Text('Ok'),
      cancellation: Empty,
      onConfirm: onConfirm,
      onCancel: onCancel,
    );
  }

  factory Confirmation.yesNo({
    Key? key,
    required Widget content,
    void Function(BuildContext)? onConfirm,
    void Function(BuildContext)? onCancel,
  }) {
    return Confirmation(
      key: key,
      content: content,
      confirmation: Text('Yes'),
      cancellation: Text('No'),
      onConfirm: onConfirm,
      onCancel: onCancel,
    );
  }

  factory Confirmation.createCancel({
    Key? key,
    required Widget content,
    void Function(BuildContext)? onConfirm,
    void Function(BuildContext)? onCancel,
  }) {
    return Confirmation(
      key: key,
      content: content,
      confirmation: Text('Create'),
      cancellation: Text('Cancel'),
      onConfirm: onConfirm,
      onCancel: onCancel,
    );
  }

  static Widget Function(Completer<void>) dangerous({
    required Widget content,
    required Future<void> Function(BuildContext) onConfirm,
  }) {
    return (completion) => Confirmation.yesNo(
      content: content,
      onConfirm: (ctx) {
        onConfirm(ctx)
            .then((_) => completion.complete())
            .catchError(completion.completeError);
      },
      onCancel: (_) => completion.complete(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    final theme = Theme.of(context);

    return _container.Container(
      padding: defaults.padding,
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.center,
        spacing: defaults.spacing,
        children: [
          content,
          Row(
            spacing: defaults.spacing,
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              InkWell(
                autofocus: true,
                onTap: onCancel == null ? null : () => onCancel!(context),
                mouseCursor: SystemMouseCursors.click,
                borderRadius: defaults.borderRadius,
                child: Container(
                  padding: defaults.padding,
                  decoration: BoxDecoration(
                    borderRadius: defaults.borderRadius,
                  ),
                  child: DefaultTextStyle(
                    style: theme.textTheme.labelLarge ?? TextStyle(),
                    child: cancellation,
                  ),
                ),
              ),
              if (confirmation != Empty)
                InkWell(
                  onTap: onConfirm == null ? null : () => onConfirm!(context),
                  mouseCursor: SystemMouseCursors.click,
                  borderRadius: defaults.borderRadius,
                  child: Container(
                    padding: defaults.padding,
                    decoration: BoxDecoration(
                      color: theme.colorScheme.primary,
                      borderRadius: defaults.borderRadius,
                    ),
                    child: DefaultTextStyle(
                      style: (theme.textTheme.labelLarge ?? TextStyle()).copyWith(
                        color: theme.colorScheme.onPrimary,
                      ),
                      child: confirmation,
                    ),
                  ),
                ),
            ],
          ),
        ],
      ),
    );
  }
}
