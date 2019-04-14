//
//  TransactionDescriptionTableViewCell.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionDescriptionTableViewCell: PHBaseTableViewCell {

    // MARK: - Private properties

    private lazy var dateLabel: UILabel = {

        let label = UILabel()
        label.translatesAutoresizingMaskIntoConstraints = false

        label.font = PHFonts.regular()
        label.textColor = PHColors.grey

        return label
    }()

    private lazy var titleLabel: UILabel = {

        let label = UILabel()
        label.translatesAutoresizingMaskIntoConstraints = false

        label.font = PHFonts.bold(ofSize: 16)
        label.textColor = PHColors.greyishBrown

        return label
    }()

    private lazy var descriptionLabel: UILabel = {

        let label = UILabel()
        label.translatesAutoresizingMaskIntoConstraints = false

        label.font = PHFonts.regular()
        label.textColor = PHColors.greyishBrown
        label.numberOfLines = 0

        return label
    }()

    // MARK: - Initialization

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)

        configure()
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Configuration

    override func configure() {
        super.configure()

        selectionStyle = .none

        addSubview(dateLabel)
        addSubview(titleLabel)
        addSubview(descriptionLabel)

        let spacing: CGFloat = 8
        let margin: CGFloat = 20
        dateLabel.topAnchor.constraint(equalTo: topAnchor, constant: margin).isActive = true
        dateLabel.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        dateLabel.rightAnchor.constraint(equalTo: rightAnchor, constant: -margin).isActive = true

        titleLabel.topAnchor.constraint(equalTo: dateLabel.bottomAnchor, constant: spacing).isActive = true
        titleLabel.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        titleLabel.rightAnchor.constraint(equalTo: rightAnchor, constant: -margin).isActive = true

        descriptionLabel.topAnchor.constraint(equalTo: titleLabel.bottomAnchor, constant: spacing).isActive = true
        descriptionLabel.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        descriptionLabel.rightAnchor.constraint(equalTo: rightAnchor, constant: -margin).isActive = true
        descriptionLabel.bottomAnchor.constraint(equalTo: bottomAnchor, constant: -margin).isActive = true
    }

    func configure(withImage
        image: UIImage?, withDate
        date: Date?, andTitle
        title: String, andDescription
        description: String) {

        imageView?.image = image
        dateLabel.text = date?.dateAndTimeString(ofStyle: .short)
        titleLabel.text = title
        descriptionLabel.text = description
    }
}
