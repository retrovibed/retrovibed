import 'dart:async';
import 'dart:io';
import 'package:flutter/material.dart';
import './screens.dart' as screens;
import "./theme.defaults.dart" as theming;

final ECONNREFUSED = 111;

class ErrorTests {
  static bool offline(Object obj) {
    return obj is SocketException && obj.osError?.errorCode == ECONNREFUSED;
  }

  static bool timeout(Object obj) {
    return obj is TimeoutException;
  }
}

class ErrorBoundary extends StatefulWidget {
  final Widget child;
  final AlignmentGeometry alignment;

  const ErrorBoundary(
    this.child, {
    super.key,
    this.alignment = Alignment.center,
  });

  static _ErrorBoundaryState? of(BuildContext context) {
    return context.findAncestorStateOfType<_ErrorBoundaryState>();
  }

  @override
  State<StatefulWidget> createState() => _ErrorBoundaryState();
}

class _ErrorBoundaryState extends State<ErrorBoundary> {
  Error? cause;
  Key _refresh = UniqueKey();

  void onError(Error err) {
    setState(() {
      cause = err;
    });
  }

  void reset() {
    setState(() {
      cause = null;
      _refresh = UniqueKey();
    });
  }

  @override
  Widget build(BuildContext context) {
    return screens.Overlay(
      key: _refresh,
      child: widget.child,
      overlay: cause,
      alignment: widget.alignment,
      onTap: cause != null ? reset : null,
    );
  }
}

class Error extends StatelessWidget {
  final Object? cause;
  final Widget child;

  const Error({super.key, required this.child, this.cause = null});

  @override
  StatelessElement createElement() {
    if (this.cause != null) {
      print("${this.cause.toString()}");
    }
    return super.createElement();
  }

  static Error text(String text) => Error(child: SelectableText(text));
  static Error unknown(Object obj) {
    return Error(
      child: SelectableText("an unexpected problem has occurred"),
      cause: obj,
    );
  }

  static Error? maybeErr(Object? obj) {
    if (obj == null) return null;
    if (obj is Error) return obj;
    return unknown(obj);
  }

  static Error offline(SocketException obj) {
    return Error(
      child: SelectableText(
        "unable to connect to daemon, is it running? check ${obj.address?.address}:${obj.port}.",
      ),
      cause: obj,
    );
  }

  static Error timeout(Object obj) {
    return Error(
      child: SelectableText(
        "timeout error: unable to complete within the expected timeframe",
      ),
      cause: obj,
    );
  }

  // pushes the error to the nearest boundary widget.
  static T Function(Object obj) boundary<T, Y>(
    BuildContext context,
    T result,
    Error Function(Y) onErr,
  ) {
    return (Object e) {
      ErrorBoundary.of(context)?.onError(onErr(e as Y));
      return result;
    };
  }

  @override
  Widget build(BuildContext context) {
    final defaults = theming.Defaults.of(context);
    return Container(
      decoration: BoxDecoration(color: defaults.danger),
      child: Center(child: child),
    );
  }
}
