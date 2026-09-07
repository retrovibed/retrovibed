import 'package:flutter/material.dart';
import '../flutterx.dart' show postframe;
import '../errors.dart' as errors;
import 'loading.dart';
import 'error.dart';

class LoadingGuard extends StatefulWidget {
  final Widget child;
  final Widget overlay;
  final BorderRadius borderRadius;

  const LoadingGuard(
    this.child, {
    super.key,
    this.overlay = Loading.Icon,
    this.borderRadius = BorderRadius.zero,
  });

  static LoadingGuardState? of(BuildContext context) {
    return context.findAncestorStateOfType<LoadingGuardState>();
  }

  @override
  State<LoadingGuard> createState() => LoadingGuardState();
}

class LoadingGuardState extends State<LoadingGuard> {
  int _count = 0;
  bool get loading => _count > 0;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void delta(int d, [String origin = 'unknown']) {
    debugPrint('LoadingGuard.delta: origin=$origin delta=$d count=$_count -> ${_count + d}');
    _count += d;
    assert(_count >= 0, 'LoadingGuard: decrement() called without a matching increment()');
    postframe(() => setState(() {}));
  }

  @override
  Widget build(BuildContext context) {
    return Loading(
      widget.child,
      loading: loading,
      maintainAnimation: false,
      maintainSize: false,
      overlay: widget.overlay,
      borderRadius: widget.borderRadius,
    );
  }
}

class LoadingBoundary extends StatefulWidget {
  final Widget child;
  final Widget cause;
  final bool loading;
  final String origin;
  const LoadingBoundary(
    this.child, {
    super.key,
    this.loading = true,
    this.cause = errors.Error.zero,
    this.origin = 'unknown',
  });

  @override
  State<LoadingBoundary> createState() => _LoadingBoundaryState();
}

class _LoadingBoundaryState extends State<LoadingBoundary> {
  LoadingGuardState? _guard;

  @override
  void initState() {
    super.initState();
    _guard = LoadingGuard.of(context);
    if (widget.loading) _guard?.delta(1, widget.origin);
  }

  @override
  void didUpdateWidget(LoadingBoundary old) {
    super.didUpdateWidget(old);
    if (old.loading == widget.loading) return;
    _guard?.delta(widget.loading ? 1 : -1, widget.origin);
  }

  @override
  void dispose() {
    if (widget.loading) _guard?.delta(-1, widget.origin);
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => ErrorScreen(widget.child, cause: widget.cause);
}
