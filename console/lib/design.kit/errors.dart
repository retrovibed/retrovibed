import 'dart:async';
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;
import './screens.dart' as screens;
import "./theme.defaults.dart" as theming;

final EDNSRESOLUTION = -2;
final ECONNREFUSED = 111;
final ENOROUTE = 113;

class ErrorTests {
  static bool offline(Object obj) {
    return obj is SocketException && (obj.osError?.errorCode == ECONNREFUSED || obj.osError?.errorCode == ENOROUTE);
  }

  static bool connectivity(Object obj) {
    if (obj is HandshakeException) return true;
    return false;
  }

  static bool dnsresolution(Object obj) {
    return obj is SocketException && (obj.osError?.errorCode == EDNSRESOLUTION);
  }

  static bool timeout(Object obj) {
    return obj is TimeoutException;
  }

  static bool websocketclosed(Object obj) {
    return obj is WebSocketException;
  }

  static bool socketclosed(Object obj) {
    return obj is SocketException && obj.message == 'Reading from a closed socket';
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

  static ErrorBoundaryState? of(BuildContext context) {
    return context.findAncestorStateOfType<ErrorBoundaryState>();
  }

  @override
  State<StatefulWidget> createState() => ErrorBoundaryState();
}

class ErrorBoundaryState extends State<ErrorBoundary> {
  Error? cause;

  void onError(Error err) {
    setState(() {
      cause = err;
    });
  }

  void reset() {
    setState(() {
      cause = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    return screens.Overlay.tappable(
      widget.child,
      overlay: cause ?? const SizedBox(),
      alignment: widget.alignment,
      onTap: cause != null ? reset : null,
    );
  }
}

class Error extends StatelessWidget {
  static const zero = const Error(child: const SizedBox(), trace: StackTrace.empty);
  final Object? cause;
  final StackTrace trace;
  final Widget child;
  final void Function()? onTap;
  final Color? color;

  const Error({
    super.key,
    required this.child,
    required this.trace,
    this.cause,
    this.onTap,
    this.color,
  });

  @override
  StatelessElement createElement() {
    final printerr = (Object? cause, StackTrace trace) {
      if (Platform.environment.containsKey('FLUTTER_TEST')) {
        return;
      }

      if (cause == null) {
        return;
      }

      debugPrint(cause.toString());
      debugPrint(trace.toString());
    };

    printerr(this.cause, this.trace);
    return super.createElement();
  }

  static Error text(String text, {StackTrace? trace, void Function()? onTap, Color? color}) =>
      Error(child: Text(text), trace: trace ?? StackTrace.current, onTap: onTap, color: color);

  static Error unknown(Object obj, {StackTrace? trace, void Function()? onTap}) {
    return Error(
      child: Text(
        "an unexpected problem has occurred",
        overflow: TextOverflow.ellipsis,
      ),
      cause: obj,
      trace: trace ?? StackTrace.current,
      onTap: onTap,
    );
  }

  static Error unauthorized(
    Object obj, {
    StackTrace? trace,
    void Function()? onTap,
    Widget? message,
    Color? color,
  }) {
    return Error(
      child: message ?? Text("you lack sufficient permissions"),
      cause: obj,
      trace: trace ?? StackTrace.current,
      onTap: onTap,
      color: color,
    );
  }

  static Error maybeErr(Object? obj, {StackTrace? trace, void Function()? onTap}) {
    if (obj == null) return Error.zero;
    if (obj is Error) return obj;
    return unknown(obj, trace: trace ?? StackTrace.current, onTap: onTap);
  }

  static Error offline(SocketException obj, {StackTrace? trace, void Function()? onTap}) {
    return Error(
      child: SelectableText(
        "unable to connect to daemon, is it running? check ${obj.address?.address}:${obj.port}.",
      ),
      cause: obj,
      trace: trace ?? StackTrace.current,
      onTap: onTap,
    );
  }

  static Error connectivity(Object obj, {StackTrace? trace, void Function()? onTap}) {
    return Error(
      child: SelectableText(
        "unable to connect, seems like there is a general connectivity issue impacting the destination.",
      ),
      cause: obj,
      trace: trace ?? StackTrace.current,
      onTap: onTap,
    );
  }

  static Error timeout(Object obj, {StackTrace? trace, void Function()? onTap}) {
    return Error(
      child: Text(
        "timeout error: unable to complete within the expected timeframe",
      ),
      cause: obj,
      trace: trace ?? StackTrace.current,
      onTap: onTap,
    );
  }

  static Error conflict(
    Object obj, {
    StackTrace? trace,
    void Function()? onTap,
    Widget? message,
    Color? color,
  }) {
    return Error(
      child: message ?? Text("a conflict occurred, the resource may already exist"),
      cause: obj,
      trace: trace ?? StackTrace.current,
      onTap: onTap,
      color: color,
    );
  }

  static Error ratelimited(
    Object obj, {
    StackTrace? trace,
    void Function()? onTap,
    Widget? message,
    Color? color,
  }) {
    return Error(
      child: message ?? Text("you're currently rate limited try again later"),
      cause: obj,
      trace: trace ?? StackTrace.current,
      onTap: onTap,
      color: color,
    );
  }

  // pushes the error to the nearest boundary widget.
  static Future<T> Function(Object obj) boundary<T, Y>(
    BuildContext context,
    T result,
    Error Function(Y, {void Function()? onTap}) onErr,
  ) {
    return (Object e) {
      final b = ErrorBoundary.of(context);
      final d = onErr(e as Y, onTap: b?.reset);
      b?.onError(d);
      return Future.value(result);
    };
  }

  void _showCauseDialog(BuildContext context) {
    showDialog(
      context: context,
      builder:
          (_) => AlertDialog(
            title: const Text('Error Details'),
            content: _ErrorDiagnostic(cause: cause, trace: trace),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text('Close'),
              ),
            ],
          ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final defaults = theming.Defaults.of(context);
    final zero = this == Error.zero;

    if (zero) return const SizedBox();

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: Listener(
        onPointerUp: onTap != null ? (_) => onTap!() : null,
        child: Container(
          alignment: Alignment.center,
          decoration: BoxDecoration(
            color: color ?? (zero ? Colors.transparent : defaults.danger),
            borderRadius: defaults.borderRadius,
          ),
          child: SelectionArea(
            child: GestureDetector(
              onLongPress: cause != null ? () => _showCauseDialog(context) : null,
              child: child,
            ),
          ),
        ),
      ),
    );
  }
}

class Errors {
  static Error httpauto(Object obj, {StackTrace? trace, void Function()? onTap}) {
    final t = trace ?? StackTrace.current;
    final code = httpx.ErrorsTest.statusCode(obj);
    // return [401, 403, 404, 409, 429, 502].contains(code);
    switch (code) {
      case 401:
      case 403:
        return Error.unauthorized(obj, trace: t, onTap: onTap);
      case 404:
        return Error.text("not found", trace: t, onTap: onTap);
      case 409:
        return Error.conflict(obj, trace: t, onTap: onTap);
      case 429:
        return Error.ratelimited(obj, trace: t, onTap: onTap);
      default:
        return Error.unknown(obj, trace: t, onTap: onTap);
    }
  }
}

class _ErrorDiagnostic extends StatelessWidget {
  final Object? cause;
  final StackTrace trace;

  const _ErrorDiagnostic({required this.cause, required this.trace});

  @override
  Widget build(BuildContext context) {
    if (cause is http.Response) {
      final r = cause as http.Response;
      return SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            SelectableText('Status: ${r.statusCode}'),
            if (r.request?.url != null) SelectableText('URL: ${r.request!.url}'),
            if (r.body.isNotEmpty) SelectableText(r.body),
            const SizedBox(height: 8),
            SelectableText(trace.toString()),
          ],
        ),
      );
    }

    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          SelectableText(cause.toString()),
          const SizedBox(height: 8),
          SelectableText(trace.toString()),
        ],
      ),
    );
  }
}
