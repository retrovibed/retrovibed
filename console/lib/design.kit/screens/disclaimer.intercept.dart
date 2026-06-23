import 'package:flutter/material.dart';
import '../modals.dart' as modals;
import 'disclaimer.dart';

// DisclaimerIntercept renders [child] behind an invisible tap-catcher until
// [cacheid] has been acknowledged. Tapping shows [overlay] (built with a
// `complete` callback) as a confirmation modal. Once the user accepts
// (`complete(true)`), the choice is persisted via [acknowledge] and the
// catcher is dropped — that gating tap is consumed, not replayed, so the
// next tap reaches [child]'s own handling directly. Once acknowledged
// (including on first build), no catcher is rendered at all and [child]
// receives taps natively.
class DisclaimerIntercept extends StatefulWidget {
  final Widget child;
  final Widget Function(void Function(bool proceed) complete) overlay;
  final String cacheid;
  final bool Function(String id) cached;
  final void Function(String id) acknowledge;

  const DisclaimerIntercept(
    this.child, {
    super.key,
    required this.cacheid,
    required this.overlay,
    this.cached = Disclaimer.disclaimerpath,
    this.acknowledge = Disclaimer.acknowledge,
  });

  @override
  State<DisclaimerIntercept> createState() => _DisclaimerInterceptState();
}

class _DisclaimerInterceptState extends State<DisclaimerIntercept> {
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  Future<void> _handleTap(BuildContext context) async {
    final proceed = await modals.asyncfn<bool>(
      context,
      (completion) => widget.overlay(completion.complete),
    );

    if (proceed != true) return;

    widget.acknowledge(widget.cacheid);
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    if (widget.cached(widget.cacheid)) {
      return widget.child;
    }

    return Stack(
      children: [
        widget.child,
        Positioned.fill(
          child: GestureDetector(
            behavior: HitTestBehavior.opaque,
            onTap: () => _handleTap(context),
          ),
        ),
      ],
    );
  }
}
