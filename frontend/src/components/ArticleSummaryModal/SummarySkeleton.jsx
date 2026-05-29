import React from 'react';

export default function SummarySkeleton() {
  return (
    <div className="flex flex-col gap-6 animate-pulse" aria-hidden="true">
      {/* Premium Gradient Header Indicator */}
      <div className="flex items-center gap-3">
        <div className="w-5 h-5 rounded-full bg-neutral-800"></div>
        <div className="h-5 w-40 bg-neutral-800 rounded-md"></div>
      </div>
      
      {/* Paragraph Skeleton block 1 */}
      <div className="space-y-3">
        <div className="h-4 w-full bg-neutral-800 rounded-md"></div>
        <div className="h-4 w-[95%] bg-neutral-800 rounded-md"></div>
        <div className="h-4 w-[90%] bg-neutral-800 rounded-md"></div>
        <div className="h-4 w-[85%] bg-neutral-800 rounded-md"></div>
      </div>

      {/* Paragraph Skeleton block 2 */}
      <div className="space-y-3 mt-4">
        <div className="h-4 w-[98%] bg-neutral-800 rounded-md"></div>
        <div className="h-4 w-[92%] bg-neutral-800 rounded-md"></div>
        <div className="h-4 w-[88%] bg-neutral-800 rounded-md"></div>
      </div>

      {/* Paragraph Skeleton block 3 */}
      <div className="space-y-3 mt-4">
        <div className="h-4 w-[96%] bg-neutral-800 rounded-md"></div>
        <div className="h-4 w-[80%] bg-neutral-800 rounded-md"></div>
      </div>
    </div>
  );
}
