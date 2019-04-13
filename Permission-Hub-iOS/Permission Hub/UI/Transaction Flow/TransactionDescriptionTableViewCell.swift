//
//  TransactionDescriptionTableViewCell.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionDescriptionTableViewCell: UITableViewCell {

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

        label.font = PHFonts.bold(ofSize: 14)
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

    private lazy var separatorView: UIView = {

        let view = UIView()
        view.translatesAutoresizingMaskIntoConstraints = false

        view.backgroundColor = PHColors.grey.withAlphaComponent(0.5)

        return view
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

    private func configure() {

        selectionStyle = .none

        addSubview(dateLabel)
        addSubview(titleLabel)
        addSubview(descriptionLabel)
        addSubview(separatorView)

        let spacing: CGFloat = 5
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

        separatorView.topAnchor.constraint(equalTo: descriptionLabel.bottomAnchor, constant: margin).isActive = true
        separatorView.leftAnchor.constraint(equalTo: leftAnchor).isActive = true
        separatorView.rightAnchor.constraint(equalTo: rightAnchor).isActive = true
        separatorView.bottomAnchor.constraint(equalTo: bottomAnchor).isActive = true
        separatorView.heightAnchor.constraint(equalToConstant: 1).isActive = true
    }

    func configure(withDate
        date: Date, andTitle
        title: String, andDescription
        description: String) {

        dateLabel.text = date.dateAndTimeString(ofStyle: .short)
        titleLabel.text = title
        descriptionLabel.text = description
    }
}
