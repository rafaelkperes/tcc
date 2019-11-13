# read data
filepath <- "./results.csv"
all <- read.csv(filepath, header = TRUE)

# filter by type
ints <- all[all$type=="int", ]
floats <- all[all$type=="float", ]
strings <- all[all$type=="strings", ]
objects <- all[all$type=="object", ]

# color sets
degrade <- c("deepskyblue4","deepskyblue3", "deepskyblue2", "deepskyblue1")
degmix <- c("deepskyblue2","deepskyblue4", "deepskyblue1", "deepskyblue3")
accessible <- c("deepskyblue4","darkgreen", "darkorchid4", "firebrick")
colorset <- degmix

# pdf settings
wid = 8
hei = 6
fontsize = 12

# labels
lTime = "Total Time (in nanoseconds)"
lTSize = "Total Size (in bytes)"
lFormat = "Data Format"

# plot
pdf(file="int-time-plot.pdf", height=hei, width=wid, pointsize = fontsize)
boxplot(dTotal~format, data=ints, col=(colorset),
        main="64-bit Integer Serialization",
        xlab=lFormat,
        ylab=lTime
)
dev.off()


